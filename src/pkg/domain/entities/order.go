package entities

import (
	"fmt"
	"time"

	"github.com/d0ze/golang-hft/src/internal"
	"github.com/shopspring/decimal"
)

type Order struct {
	Id            string               `bson:"order_id,omitempty"`
	RemoteId      int64                `bson:"remote_id,omitempty"`
	LimitPrice    decimal.Decimal      `bson:"limit_price,omitempty"`
	MarketPrice   decimal.Decimal      `bson:"market_price,omitempty"`
	Type          internal.OrderType   `bson:"order_type,omitempty"`
	PriceType     internal.PriceType   `bson:"order_type,omitempty"`
	InitialVolume decimal.Decimal      `bson:"initial_volume,omitempty"`
	Side          internal.OrderSide   `bson:"side,omitempty"`
	Status        internal.OrderStatus `bson:"order_status,omitempty"`
	Market        internal.Market      `bson:"market,omitempty"`
	CreatedAt     time.Time            `bson:"created_at,omitempty"`
	ReduceOnly    bool                 `bson:"reduceonly,omitempty"`
	Ioc           bool                 `bson:"ioc,omitempty"`
	PostOnly      bool                 `bson:"postonly,omitempty"`
	Leverage      int64                `bson:"leverage,omitempty"`
}

func (Order) TableName() string {
	return "orders"
}

func (order *Order) String() string {
	return fmt.Sprintf("%s %s %s@%s", order.Side, order.InitialVolume, order.Market, order.MarketPrice)
}

func (order *Order) IsFinal() bool {
	return order.Status != internal.CREATED && order.Status != internal.OPEN
}

func (order *Order) IsSpot() bool {
	return order.Type == internal.SPOT
}

func (order *Order) IsFuture() bool {
	return order.Type == internal.FUTURE
}

func (order *Order) IsMargin() bool {
	return order.Type == internal.MARGIN
}

func (order *Order) GetMarketCost() decimal.Decimal {
	return order.MarketPrice.Mul(order.InitialVolume)
}

func (order *Order) GetTradeCurrency() internal.Currency {
	return Markets.GetTradeCurrency(order.Market)
}

func (order *Order) GetReferenceCurrency() internal.Currency {
	return Markets.GetReferenceCurrency(order.Market)
}

func (order *Order) GetSpentCurrency() internal.Currency {
	if order.Side == internal.BUY {
		return order.GetReferenceCurrency()
	} else {
		return order.GetTradeCurrency()
	}
}
