package http

import "time"

type PlaceOrderResponse struct {
	RetCode int            `json:"retCode"`
	RetMsg  string         `json:"retMsg"`
	Result  *OrderResponse `json:"result"`
	Time    int64          `json:"time"`
}

type OrderBookResponse struct {
	RetCode int              `json:"retCode"`
	RetMsg  string           `json:"retMsg"`
	Result  *OrderBookResult `json:"result"`
	Time    int64            `json:"time"`
}

type CacelOrderResponse PlaceOrderResponse

type OrderResponse struct {
	OrderId     string `json:"orderId"`
	OrderLinkId string `json:"orderLinkId"`
}

type OrderHistoryResponse struct {
	RetCode int                `json:"retCode"`
	RetMsg  string             `json:"retMsg"`
	Result  *OrderListResponse `json:"result"`
	Time    int64              `json:"time"`
}

type BorrowHistoryResponse struct {
	RetCode int                 `json:"retCode"`
	RetMsg  string              `json:"retMsg"`
	Result  *BorrowListResponse `json:"result"`
	Time    int64               `json:"time"`
}

type OrderListResponse struct {
	NextPageCursor string   `json:"nextPageCursor"`
	Category       string   `json:"category"`
	List           []*Order `json:"list"`
}

type BorrowListResponse struct {
	NextPageCursor string    `json:"nextPageCursor"`
	Category       string    `json:"category"`
	List           []*Borrow `json:"list"`
}

type TickerListResponse struct {
	NextPageCursor string    `json:"nextPageCursor"`
	Category       string    `json:"category"`
	List           []*Ticker `json:"list"`
}

type OrderBookResult struct {
	Symbol    string     `json:"s"`
	Asks      [][]string `json:"a"`
	Bids      [][]string `json:"b"`
	Timestamp int64      `json:"ts"`
	UpdateID  int        `json:"u"`
}

type TickersResponse struct {
	RetCode int                 `json:"retCode"`
	RetMsg  string              `json:"retMsg"`
	Result  *TickerListResponse `json:"result"`
	Time    int64               `json:"time"`
}

type APIKeyInformationResponse struct {
	RetCode int                            `json:"retCode"`
	RetMsg  string                         `json:"retMsg"`
	Result  *APIKeyInformationListResponse `json:"result"`
	Time    int64                          `json:"time"`
}

type Order struct {
	Symbol             string `json:"symbol"`
	OrderType          string `json:"orderType"`
	OrderLinkID        string `json:"orderLinkId"`
	OrderID            string `json:"orderId"`
	CancelType         string `json:"cancelType"`
	AvgPrice           string `json:"avgPrice"`
	StopOrderType      string `json:"stopOrderType"`
	LastPriceOnCreated string `json:"lastPriceOnCreated"`
	OrderStatus        string `json:"orderStatus"`
	TakeProfit         string `json:"takeProfit"`
	CumExecValue       string `json:"cumExecValue"`
	TriggerDirection   int    `json:"triggerDirection"`
	BlockTradeID       string `json:"blockTradeId"`
	RejectReason       string `json:"rejectReason"`
	IsLeverage         string `json:"isLeverage"`
	Price              string `json:"price"`
	OrderIv            string `json:"orderIv"`
	CreatedTime        string `json:"createdTime"`
	TpTriggerBy        string `json:"tpTriggerBy"`
	PositionIdx        int    `json:"positionIdx"`
	TimeInForce        string `json:"timeInForce"`
	LeavesValue        string `json:"leavesValue"`
	UpdatedTime        string `json:"updatedTime"`
	Side               string `json:"side"`
	TriggerPrice       string `json:"triggerPrice"`
	CumExecFee         string `json:"cumExecFee"`
	SlTriggerBy        string `json:"slTriggerBy"`
	LeavesQty          string `json:"leavesQty"`
	CloseOnTrigger     bool   `json:"closeOnTrigger"`
	CumExecQty         string `json:"cumExecQty"`
	ReduceOnly         bool   `json:"reduceOnly"`
	Qty                string `json:"qty"`
	StopLoss           string `json:"stopLoss"`
	TriggerBy          string `json:"triggerBy"`
}
type Ticker struct {
	Symbol        string `json:"symbol"`
	Bid1Price     string `json:"bid1Price"`
	Bid1Size      string `json:"bid1Size"`
	Ask1Price     string `json:"ask1Price"`
	Ask1Size      string `json:"ask1Size"`
	LastPrice     string `json:"lastPrice"`
	PrevPrice24H  string `json:"prevPrice24h"`
	Price24HPcnt  string `json:"price24hPcnt"`
	HighPrice24H  string `json:"highPrice24h"`
	LowPrice24H   string `json:"lowPrice24h"`
	Turnover24H   string `json:"turnover24h"`
	Volume24H     string `json:"volume24h"`
	UsdIndexPrice string `json:"usdIndexPrice"`
}

type WithdrawResponse struct {
	RetCode int                 `json:"retCode"`
	RetMsg  string              `json:"retMsg"`
	Result  *WithdrawIDResponse `json:"result"`
	Time    int64               `json:"time"`
}

type WithdrawIDResponse struct {
	ID string `json:"id"`
}

type APIKeyInformationListResponse struct {
	ID            string       `json:"id"`
	Note          string       `json:"note"`
	APIKey        string       `json:"apiKey"`
	ReadOnly      int          `json:"readOnly"`
	Secret        string       `json:"secret"`
	Permissions   *Permissions `json:"permissions"`
	Ips           []string     `json:"ips"`
	Type          int          `json:"type"`
	DeadlineDay   int          `json:"deadlineDay"`
	ExpiredAt     time.Time    `json:"expiredAt"`
	CreatedAt     time.Time    `json:"createdAt"`
	Unified       int          `json:"unified"`
	Uta           int          `json:"uta"`
	UserID        int          `json:"userID"`
	InviterID     int          `json:"inviterID"`
	VipLevel      string       `json:"vipLevel"`
	MktMakerLevel string       `json:"mktMakerLevel"`
	AffiliateID   int          `json:"affiliateID"`
	RsaPublicKey  string       `json:"rsaPublicKey"`
}

type Permissions struct {
	ContractTrade []string      `json:"ContractTrade"`
	Spot          []string      `json:"Spot"`
	Wallet        []string      `json:"Wallet"`
	Options       []string      `json:"Options"`
	Derivatives   []string      `json:"Derivatives"`
	CopyTrading   []string      `json:"CopyTrading"`
	BlockTrade    []interface{} `json:"BlockTrade"`
	Exchange      []string      `json:"Exchange"`
	Nft           []string      `json:"NFT"`
}

type TransferableCoinsResponse struct {
	RetCode int                    `json:"retCode"`
	RetMsg  string                 `json:"retMsg"`
	Result  *TransferableCoinsList `json:"result"`
	Time    int64                  `json:"time"`
}

type TransferableCoinsList struct {
	List []string `json:"list"`
}

type InternalTransferResponse struct {
	RetCode int                     `json:"retCode"`
	RetMsg  string                  `json:"retMsg"`
	Result  *InternalTransferResult `json:"result"`
	Time    int64                   `json:"time"`
}

type InternalTransferResult struct {
	TransferId string `json:"transferId"`
}

type TransferWithdrawFlowResponse struct {
	WithdrawId string
	TransferId string
}

type WalletBalanceResponse struct {
	RetCode int                  `json:"retCode"`
	RetMsg  string               `json:"retMsg"`
	Result  *WalletBalanceResult `json:"result"`
	Time    int64                `json:"time"`
}

type KlineResponse struct {
	RetCode int          `json:"retCode"`
	RetMsg  string       `json:"retMsg"`
	Result  *KlineResult `json:"result"`
	Time    int64        `json:"time"`
}

type KlineResult struct {
	Symbol   string  `json:"symbol"`
	Category string  `json:"category"`
	List     []Kline `json:"list"`
}

type WalletBalanceResult struct {
	List []WalletBalance `json:"list"`
}

type WalletBalance struct {
	TotalEquity            string            `json:"totalEquity"`
	AccountIMRate          string            `json:"accountIMRate"`
	TotalMarginBalance     string            `json:"totalMarginBalance"`
	TotalInitialMargin     string            `json:"totalInitialMargin"`
	AccountType            string            `json:"accountType"`
	TotalAvailableBalance  string            `json:"totalAvailableBalance"`
	AccountMMRate          string            `json:"accountMMRate"`
	TotalPerpUPL           string            `json:"totalPerpUPL"`
	TotalWalletBalance     string            `json:"totalWalletBalance"`
	AccountLTV             string            `json:"accountLTV"`
	TotalMaintenanceMargin string            `json:"totalMaintenanceMargin"`
	Coin                   []CoinBalanceInfo `json:"coin"`
}

type CoinBalanceInfo struct {
	AvailableToBorrow   string `json:"availableToBorrow"`
	Bonus               string `json:"bonus"`
	AccruedInterest     string `json:"accruedInterest"`
	AvailableToWithdraw string `json:"availableToWithdraw"`
	TotalOrderIM        string `json:"totalOrderIM"`
	Equity              string `json:"equity"`
	TotalPositionMM     string `json:"totalPositionMM"`
	UsdValue            string `json:"usdValue"`
	UnrealisedPnl       string `json:"unrealisedPnl"`
	BorrowAmount        string `json:"borrowAmount"`
	TotalPositionIM     string `json:"totalPositionIM"`
	WalletBalance       string `json:"walletBalance"`
	CumRealisedPnl      string `json:"cumRealisedPnl"`
	Coin                string `json:"coin"`
}

type Borrow struct {
	CreatedTime               int64  `json:"createdTime"`
	CostExemption             string `json:"costExemption"`
	InterestBearingBorrowSize string `json:"InterestBearingBorrowSize"`
	Currency                  string `json:"currency"`
	HourlyBorrowRate          string `json:"hourlyBorrowRate"`
	BorrowCost                string `json:"borrowCost"`
}

type Kline []string
