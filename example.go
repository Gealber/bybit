package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Gealber/bybit/config"
	bybitHttp "github.com/Gealber/bybit/http"
	bybitWs "github.com/Gealber/bybit/websocket"
)

func main() {
	// ctx := context.Background()
	cfg := config.Config()

	// err := httpExample(ctx, cfg)
	// if err != nil {
	// 	panic(err)
	// }

	// WEBSOCKET EXAMPLE:
	// websocketExample(ctx, cfg)

	client, err := bybitHttp.New(cfg)
	if err != nil {
		panic(err)
	}

	query := bybitHttp.OrderBookParams{
		Category: "spot",
		Symbol:   "TONUSDT",
		Limit:    50,
	}
	resp, err := client.GetOrderBook(query)
	if err != nil {
		panic(err)
	}

	fmt.Println("ASKS")
	for _, r := range resp.Asks {
		fmt.Printf("RESP: %+v\n", r)
	}

	fmt.Println("BIDS")
	for _, r := range resp.Bids {
		fmt.Printf("RESP: %+v\n", r)
	}
}

func websocketExample(ctx context.Context, cfg *config.AppConfig) {
	wb := bybitWs.NewClient(cfg)

	tickerSubsciption := bybitWs.Request{
		Op: "subscribe",
		Args: []interface{}{
			// bybitWs.TickersBtcUSDTTopic,
			// bybitWs.TickersTONUSDTTopic,
			// bybitWs.TickersXLMUSDTTopic,
			bybitWs.TickersETHUSDCTopic,
			bybitWs.TickersETHUSDTTopic,
			bybitWs.TickersUSDCUSDTTopic,
		},
	}

	subscriptions := []bybitWs.Request{
		tickerSubsciption,
	}
	handlers := map[string]bybitWs.Handler{
		// bybitWs.TickersBtcUSDTTopic: bybitWs.NewTickersHandler(),
		// bybitWs.TickersTONUSDTTopic: bybitWs.NewTickersHandler(),
		// bybitWs.TickersXLMUSDTTopic: bybitWs.NewTickersHandler(),
		bybitWs.TickersETHUSDCTopic:  bybitWs.NewTickersHandler(),
		bybitWs.TickersETHUSDTTopic:  bybitWs.NewTickersHandler(),
		bybitWs.TickersUSDCUSDTTopic: bybitWs.NewTickersHandler(),
	}

	if err := wb.Run(ctx, subscriptions, handlers); err != nil {
		log.Panic(err)
	}
}

func httpExample(ctx context.Context, cfg *config.AppConfig) error {
	client, err := bybitHttp.New(cfg)
	if err != nil {
		return err
	}

	// return client.PlaceCascadeOrders(bybitHttp.BuyDirection, bybitHttp.TonChain, 0.0001, 10600)
	queryParams := bybitHttp.BorrowHistoryParams{
		Currency: "USDT",
	}
	list, err := client.BorrowHistory(queryParams)
	if err != nil {
		return err
	}

	for _, record := range list {
		log.Printf("RECORD: %+v\n", record)
	}

	return nil
}

func cancelOrder(ctx context.Context, cfg *config.AppConfig) error {
	client, err := bybitHttp.New(cfg)
	if err != nil {
		return err
	}

	orders := []bybitHttp.OrderResponse{}

	for _, order := range orders {
		cancel := bybitHttp.CancelRequest{
			Category:    "spot",
			Symbol:      "TONUSDT",
			OrderID:     order.OrderId,
			OrderLinkId: order.OrderLinkId,
		}
		_, err := client.CancelOrder(cancel)
		if err != nil {
			// return err
			log.Println(err)
		}
	}

	return nil
}
