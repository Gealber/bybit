package websocket

type ChannelType string
type CoverType string

const (
	PublicChannel  ChannelType = "public"
	PrivateChannel ChannelType = "private"
	Spot           CoverType   = "spot"
	Linear         CoverType   = "linear"
	Inverse        CoverType   = "inverse"
	Option         CoverType   = "option"
)

const (
	TickersTONUSDTTopic = "tickers.TONUSDT"
	TickersBtcUSDTTopic = "tickers.BTCUSDT"
)

const (
	APIVersion           = "v5"
	ByBitWebsocketDomain = "stream.bybit.com"
	ByBitSpotPath        = "/spot"
)

const (
	PingTimeout           = 20
	TickerKeyTimeout      = 45
	MaxRetrialConnections = 10
)
