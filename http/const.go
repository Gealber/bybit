package http

const (
	RecvWindow = "6000"
	BaseURL    = "https://api.bybit.com"
	APIVersion = "v5"
	RetCodeOK  = 0
)

const (
	SpotCategory = "spot"

	// direction of order.
	BuyDirection  = "Buy"
	SellDirection = "Sell"

	// order type.
	MarketOrder = "Market"
	LimitOrder  = "Limit"

	// account types.
	UnifiedAccount = "UNIFIED"
	FundingAccount = "FUNDING"

	TonChain      = "TON"
	TonUSDTSymbol = "TONUSDT"
)

const (
	DeafaultPlaceOrdersQty = 10
)