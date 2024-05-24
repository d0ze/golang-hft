package entities

import "github.com/shopspring/decimal"

type Balance struct {
	TradeBalance  decimal.Decimal // combined balance of all equity currencies
	InitialMargin decimal.Decimal // margin amount used for open positions
	FreeMargin    decimal.Decimal // available margin for new operations (equity - margin)
	Equity        decimal.Decimal // trade balance + unrealized P&L
	MarginLevel   decimal.Decimal // margin level of the account (equity / margin) * 100
}
