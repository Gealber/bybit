package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/Gealber/bybit/config"

	"github.com/gorilla/websocket"
)

// Client represents connection with ByBit Websocket API.
type Client struct {
	APIKey    string
	APISecret string
	logger    *log.Logger
}

// Handler for processing message
type Handler interface {
	ProcessMsg(ctx context.Context, obj any) error
}

// NewClient creates a new websocket client
func NewClient(cfg *config.AppConfig) *Client {
	// bybit WS logger.
	bybitWSLogger := log.New(os.Stdout, "[bybit-ws]", log.Lshortfile)

	return &Client{
		APIKey:    cfg.ByBit.APIKey,
		APISecret: cfg.ByBit.APISecret,
		logger:    bybitWSLogger,
	}
}

func (c *Client) path(channelType ChannelType, operation CoverType) string {
	return fmt.Sprintf("/%s/%s/%s", APIVersion, channelType, operation)
}

func (c *Client) connect() (*websocket.Conn, error) {
	spotPath := c.path(PublicChannel, Spot)

	u := url.URL{Scheme: "wss", Host: ByBitWebsocketDomain, Path: spotPath}
	c.logger.Printf("connecting to %s", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	return conn, err
}

func sendPing(conn *websocket.Conn) error {
	pingReq := Request{
		ReqID: "100001",
		Op:    "ping",
	}

	return conn.WriteJSON(&pingReq)
}

// Run connect to bybit websocket, general idea of what it does.
// 1. Subscribe to tickers
// 2. Read message from websocket.
// 3. Send every 20 seconds a ping, to avoid disconnections.
// 4. In case of abnormal close of connection, performs a reconnection.
// 5. In case the reconnection exceed the max allowed, shut the program.
// 6. Also listen to Ctr+C commands to shutdown gratefully.
func (c *Client) Run(
	ctx context.Context,
	subscriptions []Request,
	handlers map[string]Handler,
) error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	pingTicker := time.NewTicker(PingTimeout * time.Second)
	errChn := make(chan error)
	done := make(chan struct{})
	var connections int = 0

CONNECTION:
	conn, err := c.connect()
	if err != nil {
		return err
	}
	defer func() {
		pingTicker.Stop()
		close(errChn)
		conn.Close()
	}()

	connections++

	go func() {
		defer close(done)
		err := c.processRead(ctx, done, conn, subscriptions, handlers)
		if err != nil {
			errChn <- err

			return
		}
	}()

	// wait for possible errors and interrupt signals.
	// send ping command every PingTimeout seconds.
	for {
		select {
		case err := <-errChn:
			if waitTime, ok := retriableError(err, connections); ok {
				if connections > MaxRetrialConnections {
					return err
				}

				c.logger.Println("RECONNECTING....ABNORMAL CLOUSURE....")
				c.logger.Printf("Sleeping %v milliseconds\n", waitTime)
				time.Sleep(waitTime)
				goto CONNECTION
			}

			return err
		case <-interrupt:
			return c.handleInterruptSignal(done, conn)
		case <-pingTicker.C:
			err := sendPing(conn)
			if err != nil {
				return err
			}
		}
	}
}

func (c *Client) processMsg(ctx context.Context, data []byte, handlers map[string]Handler) error {
	var msg PublicResponse

	err := json.Unmarshal(data, &msg)
	if err != nil {
		return err
	}

	switch msg.Topic {
	case TickersTONUSDTTopic:
		tickersHandler := handlers[TickersTONUSDTTopic]
		err := processTickerTopic(ctx, data, tickersHandler)
		if err != nil {
			return err
		}
	}

	c.logger.Printf("MSG: %s\n", string(data))
	return nil
}

func (c *Client) processRead(
	ctx context.Context,
	done chan struct{},
	conn *websocket.Conn,
	subscriptions []Request,
	handlers map[string]Handler,
) error {
	// first ping to send.
	sendPing(conn)

	for _, subscription := range subscriptions {
		err := conn.WriteJSON(subscription)
		if err != nil {
			return fmt.Errorf("sending subscription %w", err)
		}
	}

	for {
		select {
		case <-done:
			return nil
		default:
			_, message, err := conn.ReadMessage()
			if err != nil {
				return fmt.Errorf("reading %w", err)
			}

			err = c.processMsg(ctx, message, handlers)
			if err != nil {
				c.logger.Println("ERR PROCESSING READ: ", err.Error())
			}
		}
	}
}

func (c *Client) handleInterruptSignal(
	done chan struct{},
	conn *websocket.Conn,
) error {
	c.logger.Println("Closing connection...it might take a few seconds")
	done <- struct{}{}
	// Cleanly close the connection by sending a close message and then
	// waiting (with timeout) for the server to close the connection.
	err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		return err
	}

	select {
	case <-time.After(time.Second):
	}

	return nil
}

func processTickerTopic(ctx context.Context, data []byte, handler Handler) error {
	var tickersMsg TickersResponse

	err := json.Unmarshal(data, &tickersMsg)
	if err != nil {
		return err
	}

	return handler.ProcessMsg(ctx, tickersMsg)
}

func retriableError(err error, connections int) (time.Duration, bool) {
	waitTime := time.Duration(connections) * 500 * time.Millisecond
	// 1006 is a reserved value and MUST NOT be set as a status code in a
	// Close control frame by an endpoint.  It is designated for use in
	// applications expecting a status code to indicate that the
	// connection was closed abnormally, e.g., without sending or
	// receiving a Close control frame
	isAbnormalClosure := websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure)
	// 1001 indicates that an endpoint is "going away", such as a server
	// going down or a browser having navigated away from a page.
	isGoingAway := websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway)

	isRetriable := isAbnormalClosure || isGoingAway

	if websocket.IsUnexpectedCloseError(err, websocket.CloseTryAgainLater) {
		isRetriable = true
		waitTime = time.Duration(connections) * time.Second
	}

	return waitTime, isRetriable
}
