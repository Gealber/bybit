package http

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Gealber/bybit/config"
	query "github.com/google/go-querystring/query"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

// Client represents connection with ByBit REST API.
type Client struct {
	APIKey     string
	APISecret  string
	HttpClient http.Client
	BaseURL    string
	logger     *log.Logger
}

// New create a new instance of a client
func New(cfg *config.AppConfig) (*Client, error) {
	httpClient := http.Client{}

	bybitLoggerHTTP := log.New(os.Stdout, "[bybit-http]", log.Lshortfile)

	if cfg.ByBit.APIKey == "" {
		return nil, errors.New("empty api key in env checkout environment variable BYBIT_API_KEY")
	}

	if cfg.ByBit.APISecret == "" {
		return nil, errors.New("empty api secret in env checkout environment variable BYBIT_API_SECRET")
	}

	return &Client{
		APIKey:     cfg.ByBit.APIKey,
		APISecret:  cfg.ByBit.APISecret,
		HttpClient: httpClient,
		BaseURL:    cfg.ByBit.BaseURL,
		logger:     bybitLoggerHTTP,
	}, nil
}

// Do performs http request according to the req provided
// the response is stored in the pointer to a struct 'objResp'
func (c *Client) Do(
	req *http.Request,
	objResp interface{},
) error {
	response, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return ErrorUnexpectedStatus
	}

	b, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, objResp)
}

// NewRequest creates a new request with the arguments provided
// objBody should be a pointer to struct with json tags
// this param represent the body to be sent in a POST request
// GET: only method, path, and queryParams
func (c *Client) NewRequest(
	method, path string,
	queryParams any,
	objBody any,
) (*http.Request, error) {
	var (
		err        error
		bodyReader io.Reader
		queries    url.Values
	)

	if queryParams != nil {
		queries, err = query.Values(queryParams)
		if err != nil {
			return nil, err
		}
	}

	timestamp := time.Now().UTC().UnixMilli()
	sign := c.genSignHash(timestamp, queries.Encode())
	url := c.buildURL(path, queries)

	c.logger.Println("URL: ", url)

	if objBody != nil {
		data, err := json.Marshal(objBody)
		if err != nil {
			return nil, err
		}

		sign = c.genSignHash(timestamp, string(data))
		bodyReader = bytes.NewBuffer(data)
	}

	request, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	request.Header.Add("X-BAPI-API-KEY", c.APIKey)
	request.Header.Add("X-BAPI-TIMESTAMP", fmt.Sprintf("%d", timestamp))
	request.Header.Add("X-BAPI-SIGN", sign)
	request.Header.Add("X-BAPI-RECV-WINDOW", RecvWindow)

	return request, nil
}

// PlaceOrder place an order in the exchange
func (c *Client) PlaceOrder(order OrderRequest) (*OrderResponse, error) {
	path := "order/create"

	request, err := c.NewRequest(http.MethodPost, path, nil, &order)
	if err != nil {
		return nil, err
	}

	var response PlaceOrderResponse

	err = c.Do(request, &response)
	if err != nil {
		return nil, err
	}

	if response.RetCode != RetCodeOK {
		return nil, errors.New(response.RetMsg)
	}

	return response.Result, nil
}

// CancelOrder cancel an order in the exchange
func (c *Client) CancelOrder(cancel CancelRequest) (*OrderResponse, error) {
	path := "order/cancel"

	request, err := c.NewRequest(http.MethodPost, path, nil, &cancel)
	if err != nil {
		return nil, err
	}

	var response CacelOrderResponse

	err = c.Do(request, &response)
	if err != nil {
		return nil, err
	}

	if response.RetCode != RetCodeOK {
		return nil, errors.New(response.RetMsg)
	}

	return response.Result, nil
}

// OrderHistory retrieve the order history
func (c *Client) OrderHistory(queryParams HistoryParams) ([]*Order, error) {
	path := "order/history"

	request, err := c.NewRequest(http.MethodGet, path, queryParams, nil)
	if err != nil {
		return nil, err
	}

	var response OrderHistoryResponse
	err = c.Do(request, &response)
	if err != nil {
		return nil, err
	}

	if response.RetCode != RetCodeOK {
		return nil, errors.New(response.RetMsg)
	}

	return response.Result.List, nil
}

// OpenOrders retreive open orders
func (c *Client) OpenOrders(queryParams any) ([]*Order, error) {
	path := "order/realtime"

	request, err := c.NewRequest(http.MethodGet, path, queryParams, nil)
	if err != nil {
		return nil, err
	}

	var response OrderHistoryResponse
	err = c.Do(request, &response)
	if err != nil {
		return nil, err
	}

	if response.RetCode != RetCodeOK {
		return nil, errors.New(response.RetMsg)
	}

	return response.Result.List, nil
}

// GetTickers retrieve tickers of a given symbol specified in queryParams
func (c *Client) GetTickers(queryParams TickerParams) ([]*Ticker, error) {
	path := "market/tickers"

	request, err := c.NewRequest(http.MethodGet, path, queryParams, nil)
	if err != nil {
		return nil, err
	}

	var response TickersResponse
	err = c.Do(request, &response)
	if err != nil {
		return nil, err
	}

	return response.Result.List, err
}

// Witdraw create a withdraw request. Take into account that to perform a withdraw
// you api key should be bind to a fix IP address. Read more about in bybit doc.
func (c *Client) Withdraw(withdraw WithdrawRequest) (string, error) {
	path := "asset/withdraw/create"

	request, err := c.NewRequest(http.MethodPost, path, nil, &withdraw)
	if err != nil {
		return "", err
	}

	var response WithdrawResponse
	err = c.Do(request, &response)
	if err != nil {
		return "", err
	}

	if response.RetCode != RetCodeOK {
		return "", errors.New(response.RetMsg)
	}

	return response.Result.ID, err
}

// GetAPIKeyINformation retrieve api key information
func (c *Client) GetAPIKeyInformation() (*APIKeyInformationListResponse, error) {
	path := "user/query-api"

	request, err := c.NewRequest(http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}

	var response APIKeyInformationResponse
	err = c.Do(request, &response)
	if err != nil {
		return nil, err
	}

	return response.Result, err
}

// GetTransferableCoins retreive transferable coins
func (c *Client) GetTransferableCoins(query TransferableCoinsListParams) (*TransferableCoinsList, error) {
	path := "asset/transfer/query-transfer-coin-list"

	request, err := c.NewRequest(http.MethodGet, path, &query, nil)
	if err != nil {
		return nil, err
	}

	var response TransferableCoinsResponse
	err = c.Do(request, &response)
	if err != nil {
		return nil, err
	}

	return response.Result, err
}

// CreateInternalTransfer an internal transfer
func (c *Client) CreateInternalTransfer(transfer TransferRequest) (string, error) {
	path := "asset/transfer/inter-transfer"

	request, err := c.NewRequest(http.MethodPost, path, nil, &transfer)
	if err != nil {
		return "", err
	}

	var response InternalTransferResponse
	err = c.Do(request, &response)
	if err != nil {
		return "", err
	}

	return response.Result.TransferId, err
}

// TransferWithdrawFlow make an internal transfer from Unified to Funding account
// and then withdraw the that token
func (c *Client) TransferWithdrawFlow(coin, chain, amount, address string) (*TransferWithdrawFlowResponse, error) {
	transfer := TransferRequest{
		TransferID:      uuid.New().String(),
		Coin:            coin,
		Amount:          amount,
		FromAccountType: UnifiedAccount,
		ToAccountType:   FundingAccount,
	}

	transferID, err := c.CreateInternalTransfer(transfer)
	if err != nil {
		return nil, err
	}

	withdraw := WithdrawRequest{
		Coin:    coin,
		Chain:   chain,
		Address: address,
		Amount:  amount, // need to be changed
	}

	withdrawID, err := c.Withdraw(withdraw)
	if err != nil {
		return nil, err
	}

	return &TransferWithdrawFlowResponse{
		WithdrawId: withdrawID,
		TransferId: transferID,
	}, nil
}

// GetWalletBalance retrieve wallet balance
func (c *Client) GetWalletBalance(queryParams WalletBalanceParams) (*WalletBalanceResult, error) {
	path := "account/wallet-balance"

	request, err := c.NewRequest(http.MethodGet, path, queryParams, nil)
	if err != nil {
		return nil, err
	}

	var response WalletBalanceResponse
	err = c.Do(request, &response)
	if err != nil {
		return nil, err
	}

	return response.Result, err
}

// PlaceCascadeOrders is a custom method to perform several orders
// In case we want to SELL the orders will increase in value from the first bid in the order book
// In case we want to BUY the orders will decrease in value from the first ask in the order book
// The number of orders to be created is 10. All of the order created are limit orders
func (c *Client) PlaceCascadeOrders(side, coin string, priceStep, coinQty float64) error {
	currentPrice, err := c.getLatestOrderBookPrice(side, coin)
	if err != nil {
		return err
	}

	coinEquity, usdtEquity, err := c.getCoinUSDTEquity(coin)
	if err != nil {
		return err
	}

	if side == SellDirection && coinEquity < coinQty {
		return fmt.Errorf("%v: coint equity: %f quantity: %f", ErrorInsuficcientBalance, coinEquity, coinQty)
	}

	usdtToSpendBuying := func() float64 {
		spenditure := 0.0
		nextPrice := currentPrice
		for i := 0; i < DeafaultPlaceOrdersQty; i++ {
			spenditure += nextPrice * coinQty / DeafaultPlaceOrdersQty
			nextPrice += priceStep
		}

		return spenditure
	}()

	if side == BuyDirection && usdtEquity < usdtToSpendBuying {
		return fmt.Errorf("%v: usdt equity: %f quantiy x price: %f price: %f quantity: %f", ErrorInsuficcientBalance, usdtEquity, coinQty*currentPrice, currentPrice, coinQty)
	}

	orders := c.prepareCascadeOrders(side, coin, coinQty, currentPrice, priceStep)

	// perform cascade orders in goroutines.
	errsGroup, _ := errgroup.WithContext(context.Background())
	for _, order := range orders {
		orderReq := order
		errsGroup.Go(func() error {
			resp, err := c.PlaceOrder(orderReq)
			if err != nil {
				return err
			}

			c.logger.Printf("ORDER: %+v\n", resp)

			return nil
		})
	}

	return errsGroup.Wait()
}

func (c *Client) prepareCascadeOrders(side, coin string, quantity, startPrice, priceStep float64) []OrderRequest {
	remaining := quantity
	orderSize := quantity / DeafaultPlaceOrdersQty
	price := startPrice
	orders := make([]OrderRequest, 0)
	for remaining > 0 {
		orderLinkID := fmt.Sprintf("%s-%s-%s", coin, side, uuid.New().String())
		orderRequest := OrderRequest{
			Category:    SpotCategory,
			Side:        side,
			Symbol:      fmt.Sprintf("%sUSDT", strings.ToUpper(coin)),
			OrderType:   LimitOrder,
			OrderLinkId: orderLinkID,
			Qty:         strconv.FormatFloat(orderSize, 'f', 4, 64),
			Price:       strconv.FormatFloat(price, 'f', 4, 64),
		}
		if remaining < orderSize {
			orderRequest.Qty = strconv.FormatFloat(remaining, 'f', 4, 64)
			orders = append(orders, orderRequest)
			break
		}

		orders = append(orders, orderRequest)
		remaining -= orderSize
		if side == BuyDirection {
			price -= priceStep
		} else {
			price += priceStep
		}
	}

	return orders
}

func (c *Client) getCoinUSDTEquity(coin string) (float64, float64, error) {
	queryParams := WalletBalanceParams{
		AccountType: UnifiedAccount,
		Coin:        coin,
	}

	balanceInfo, err := c.GetWalletBalance(queryParams)
	if err != nil {
		return 0, 0, err
	}

	var (
		coinEquity float64
		usdtEquity float64
	)

	for _, balance := range balanceInfo.List {
		if len(balance.Coin) == 0 {
			return 0, 0, ErrorUnavailableInformation
		}

		usdtEquity, err = strconv.ParseFloat(balance.TotalAvailableBalance, 64)
		if err != nil {
			return 0, 0, err
		}

		for _, coinInfo := range balance.Coin {
			coinEquity, err = strconv.ParseFloat(coinInfo.Equity, 64)
			if err != nil {
				return 0, 0, err
			}
		}
	}

	return coinEquity, usdtEquity, nil
}

func (c *Client) getLatestOrderBookPrice(side, coin string) (float64, error) {
	tickersParams := TickerParams{
		Category: SpotCategory,
		Symbol:   fmt.Sprintf("%sUSDT", coin),
	}

	tickers, err := c.GetTickers(tickersParams)
	if err != nil {
		return 0, err
	}

	var currentPrice float64
	for _, ticker := range tickers {
		switch side {
		case SellDirection:
			bidPrice, err := strconv.ParseFloat(ticker.Bid1Price, 64)
			if err != nil {
				return 0, err
			}

			currentPrice = bidPrice
		case BuyDirection:
			askPrice, err := strconv.ParseFloat(ticker.Ask1Price, 64)
			if err != nil {
				return 0, err
			}

			currentPrice = askPrice
		}
	}

	return currentPrice, nil
}

func (c *Client) genSignHash(timestamp int64, payload string) string {
	h := hmac.New(sha256.New, []byte(c.APISecret))

	paramStr := fmt.Sprintf("%d%s%s%s", timestamp, c.APIKey, RecvWindow, payload)

	h.Write([]byte(paramStr))

	return hex.EncodeToString(h.Sum(nil))
}

func (c *Client) buildURL(path string, queryValues url.Values) string {
	urlPath := fmt.Sprintf("%s/%s/%s", c.BaseURL, APIVersion, path)

	return fmt.Sprintf("%s?%s", urlPath, queryValues.Encode())
}
