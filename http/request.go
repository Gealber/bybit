package http

// OrderRequest entity for creating an order
type OrderRequest struct {
	Category         string `json:"category,omitempty"`
	Symbol           string `json:"symbol,omitempty"`
	IsLeverage       string `json:"isLeverage,omitempty"`
	Side             string `json:"side,omitempty"`
	OrderType        string `json:"orderType,omitempty"`
	Qty              string `json:"qty,omitempty"`
	Price            string `json:"price,omitempty"`
	TriggerDirection int    `json:"triggerDirection,omitempty"`
	OrderFilter      string `json:"orderFilter,omitempty"`
	TriggerPrice     string `json:"triggerPrice,omitempty"`
	TriggerBy        string `json:"triggerBy,omitempty"`
	OrderIv          string `json:"orderIv,omitempty"`
	TimeInForce      string `json:"timeInForce,omitempty"`
	PositionIdx      int    `json:"positionIdx,omitempty"`
	OrderLinkId      string `json:"orderLinkId,omitempty"`
	TakeProfit       string `json:"takeProfit,omitempty"`
	StopLoss         string `json:"stopLoss,omitempty"`
	TpTriggerBy      string `json:"tpTriggerBy,omitempty"`
	SlTriggerBy      string `json:"slTriggerBy,omitempty"`
	ReduceOnly       bool   `json:"reduceOnly,omitempty"`
	CloseOnTrigger   bool   `json:"closeOnTrigger,omitempty"`
	MMP              bool   `json:"mmp,omitempty"`
}

// CancelRequest entity for cancelling order
type CancelRequest struct {
	Category    string `json:"category"`
	Symbol      string `json:"symbol"`
	OrderID     string `json:"orderId,omitempty"`
	OrderLinkId string `json:"orderLinkId,omitempty"`
	OrderFilter string `json:"orderFilter,omitempty"`
}

// HistoryParams entitity for requesting history of a transaction
// used in OrderHistory
type HistoryParams struct {
	Category    string `url:"category,omitempty"`
	Symbol      string `url:"symbol,omitempty"`
	BaseCoin    string `url:"baseCoin,omitempty"`
	OrderId     string `url:"orderId,omitempty"`
	OrderLinkId string `url:"orderLinkId,omitempty"`
	OrderFilter string `url:"orderFilter,omitempty"`
	OrderStatus string `url:"orderStatus,omitempty"`
	StartTime   int64  `url:"startTime,omitempty"`
	EndTime     int64  `url:"endTime,omitempty"`
	Limit       int    `url:"limit,omitempty"`
	Cursor      string `url:"cursor,omitempty"`
}

// BorrowHistoryParams entitity for requesting history of a borrowed actions
type BorrowHistoryParams struct {
	Currency  string `url:"currency,omitempty"`
	StartTime int64  `url:"startTime,omitempty"`
	EndTime   int64  `url:"endTime,omitempty"`
	Limit     int    `url:"limit,omitempty"`
	Cursor    string `url:"cursor,omitempty"`
}

// TickerParams entity for requesting ticker information about a coin
type TickerParams struct {
	Category string `url:"category,omitempty"`
	Symbol   string `url:"symbol,omitempty"`
	BaseCoin string `url:"baseCoin,omitempty"`
	ExpDate  string `url:"expDate,omitempty"`
}

// KlineParams entity for requesting kline information about a coin
type KlineParams struct {
	Category string `url:"category"`
	Symbol   string `url:"symbol"`
	Interval string `url:"interval"`
	Start    int    `url:"start"`
	End      int    `url:"end,omitempty"`
	Limit    int    `url:"limit,omitempty"`
}

// OrderBookParams entity for requesting order book information about a coin
type OrderBookParams struct {
	Category string `url:"category"`
	Symbol   string `url:"symbol"`
	Limit    int    `url:"limit,omitempty"`
}

// WithdrawRequest entity for withdrawing assets
type WithdrawRequest struct {
	Coin        string `json:"coin"`
	Chain       string `json:"chain"`
	Address     string `json:"address"`
	Tag         string `json:"tag"`
	Amount      string `json:"amount" `
	Timestamp   int64  `json:"timestamp"`
	ForceChain  int    `json:"forceChain"`
	AccountType string `json:"accountType"`
}

type TransferableCoinsListParams struct {
	FromAccountType string `url:"fromAccountType"`
	ToAccountType   string `url:"toAccountType"`
}

type TransferRequest struct {
	TransferID      string `json:"transferId"`
	Coin            string `json:"coin"`
	Amount          string `json:"amount"`
	FromAccountType string `json:"fromAccountType"`
	ToAccountType   string `json:"toAccountType"`
}

type WalletBalanceParams struct {
	AccountType string `url:"accountType"`
	Coin        string `url:"coin"`
}
