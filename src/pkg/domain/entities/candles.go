package entities

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type Candle struct {
	Open      decimal.Decimal
	High      decimal.Decimal
	Low       decimal.Decimal
	Close     decimal.Decimal
	Timestamp time.Time
}

func NewCandle(o, h, l, c decimal.Decimal, ts time.Time) Candle {
	return Candle{
		Open:      o,
		High:      h,
		Low:       l,
		Close:     c,
		Timestamp: ts,
	}
}

func (c *Candle) String() string {
	return fmt.Sprintf("(%s) open %s; high %s; low %s; close %s",
		c.Timestamp.String(),
		c.Open.StringFixed(2),
		c.High.StringFixed(2),
		c.Low.StringFixed(2),
		c.Close.StringFixed(2))
}

func (c *Candle) IsUp() bool {
	return c.Close.GreaterThan(c.Open)
}

func (c *Candle) IsDown() bool {
	return c.Close.LessThan(c.Open)
}

// returns the meaningful price for the candle
// meaningful price is computed as o + h + l + c / 4
func (c *Candle) GetPrice() decimal.Decimal {
	return (c.High.Add(c.Low).Add(c.Open).Add(c.Close)).Div(decimal.NewFromInt(4))
}
