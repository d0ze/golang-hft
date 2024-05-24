package entities

import (
	"fmt"

	"github.com/d0ze/golang-hft/src/internal"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

// this module handles markets metadata from each exchange
// it exposes methods to retrieve (for each market defined in the application) the following informations
//   - trade currency
//   - reference currency
//   - decimals (precision)
//   - min cost of an order
type metadata struct {
	Decimals          int
	TradeCurrency     internal.Currency
	ReferenceCurrency internal.Currency
	MinCost           decimal.Decimal
	OrderMin          decimal.Decimal
}

type IMarkets interface {
	GetDecimals(market internal.Market) int
	SetDecimals(market internal.Market, value int)
	GetTradeCurrency(market internal.Market) internal.Currency
	SetTradeCurrency(market internal.Market, value internal.Currency)
	GetReferenceCurrency(market internal.Market) internal.Currency
	SetReferenceCurrency(market internal.Market, value internal.Currency)
	GetMetadata(market internal.Market) *metadata
	GetMinCost(market internal.Market) decimal.Decimal
	SetMinCost(market internal.Market, value decimal.Decimal)
	GetOrderMin(market internal.Market) decimal.Decimal
	SetOrderMin(market internal.Market, value decimal.Decimal)
	SetMetadata(
		market internal.Market,
		decimals int,
		tradeCurrency internal.Currency,
		referenceCurrency internal.Currency,
		minCost decimal.Decimal,
		orderMin decimal.Decimal)
	String() string
}

type markets struct {
	markets map[internal.Market]*metadata
}

func NewMarkets() IMarkets {
	return &markets{
		markets: map[internal.Market]*metadata{
			internal.XBTEUR:  {},
			internal.XBTUSD:  {},
			internal.XBTUSDT: {},
			internal.ETHEUR:  {},
			internal.ETHUSD:  {},
		},
	}
}

func (c *markets) GetDecimals(market internal.Market) int {
	return c.markets[market].Decimals
}

func (c *markets) SetDecimals(market internal.Market, value int) {
	if v, ok := c.markets[market]; ok {
		v.Decimals = value
	}
}

func (c *markets) GetTradeCurrency(market internal.Market) internal.Currency {
	return c.markets[market].TradeCurrency
}

func (c *markets) SetTradeCurrency(market internal.Market, value internal.Currency) {
	if v, ok := c.markets[market]; ok {
		v.TradeCurrency = value
	}
}

func (c *markets) GetReferenceCurrency(market internal.Market) internal.Currency {
	return c.markets[market].ReferenceCurrency
}

func (c *markets) SetReferenceCurrency(market internal.Market, value internal.Currency) {
	if v, ok := c.markets[market]; ok {
		v.ReferenceCurrency = value
	}
}

func (c *markets) GetMetadata(market internal.Market) *metadata {
	return c.markets[market]
}

func (c *markets) GetMinCost(market internal.Market) decimal.Decimal {
	return c.markets[market].MinCost
}

func (c *markets) SetMinCost(market internal.Market, value decimal.Decimal) {
	if v, ok := c.markets[market]; ok {
		v.MinCost = value
	}
}

func (c *markets) GetOrderMin(market internal.Market) decimal.Decimal {
	return c.markets[market].OrderMin
}

func (c *markets) SetOrderMin(market internal.Market, value decimal.Decimal) {
	if v, ok := c.markets[market]; ok {
		v.OrderMin = value
	}
}

func (c *markets) SetMetadata(
	market internal.Market,
	precision int,
	tradeCurrency internal.Currency,
	referenceCurrency internal.Currency,
	minCost decimal.Decimal,
	orderMin decimal.Decimal) {
	if v, ok := c.markets[market]; ok {
		logrus.Infof("setting metadata %d %s %s %s %s %s", precision, market, tradeCurrency, referenceCurrency, minCost, orderMin)
		v.Decimals = precision
		v.TradeCurrency = tradeCurrency
		v.ReferenceCurrency = referenceCurrency
		v.MinCost = minCost
		v.OrderMin = orderMin
	} else {
		new := &metadata{
			Decimals:          precision,
			TradeCurrency:     tradeCurrency,
			ReferenceCurrency: referenceCurrency,
			MinCost:           minCost,
			OrderMin:          orderMin,
		}
		c.markets[market] = new
	}
}

func (c *markets) String() string {
	var res string
	for market, md := range c.markets {
		res += fmt.Sprintf("[%s]: precision %d, tc %s, rc %s, mc %s", market, md.Decimals, md.TradeCurrency, md.ReferenceCurrency, md.MinCost)
	}
	return res
}

var Markets IMarkets
