package entities

import (
	"fmt"
	"time"

	"github.com/d0ze/golang-hft/src/internal"
	"github.com/shopspring/decimal"
)

type Position struct {
	Id         string
	Size       decimal.Decimal
	Side       internal.OrderSide
	OpenPrice  decimal.Decimal
	ClosePrice decimal.Decimal
	Market     internal.Market
	Realized   decimal.Decimal
	Status     internal.PositionStatus
	CreatedAt  time.Time
	Cost       decimal.Decimal
	Leverage   int
}

func (p *Position) String() string {
	return fmt.Sprintf("%s %s %s @%s", p.Side, p.Market, p.Size, p.OpenPrice)
}
