package internal

type OrderSide string
type OrderStatus string
type OrderType string
type PriceType string
type PositionStatus string
type PositionSide string
type Currency string
type Market string
type Strategy string

const (
	CREATED   OrderStatus = "created"
	OPEN      OrderStatus = "open"
	FILLED    OrderStatus = "filled"
	CANCELLED OrderStatus = "cancelled"
	ERROR     OrderStatus = "error"
)

const (
	SPOT   OrderType = "spot"
	MARGIN OrderType = "margin"
	FUTURE OrderType = "future"
)

const (
	MARKET            PriceType = "market"
	LIMIT             PriceType = "limit"
	TAKE_PROFIT       PriceType = "take-profit"
	TAKE_PROFIT_LIMIT PriceType = "take-profit-limit"
	STOP_LOSS         PriceType = "stop-loss"
	STOP_LOSS_LIMIT   PriceType = "stop-loss-limit"
)

const (
	BUY  OrderSide = "buy"
	SELL OrderSide = "sell"
)

const (
	POPEN  PositionStatus = "open"
	PCLOSE PositionStatus = "closed"
)

const (
	LONG  PositionSide   = "long"
	SHORT PositionStatus = "short"
)

const (
	XBT  Currency = "XBT"
	EUR  Currency = "EUR"
	USD  Currency = "USD"
	ETH  Currency = "ETH"
	USDT Currency = "USDT"
	LTC  Currency = "LTC"
)
const (
	XBTEUR  Market = "XBTEUR"
	XBTUSD  Market = "XBTUSD"
	XBTUSDT Market = "XBTUSDT"
	ETHEUR  Market = "ETHEUR"
	ETHUSD  Market = "ETHUSD"
	LTCEUR  Market = "LTCEUR"
)

const (
	TWAP    Strategy = "twap"
	CHATGPT Strategy = "cgpt"
)
