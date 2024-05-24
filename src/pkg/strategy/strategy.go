package strategy

import (
	"time"

	"github.com/d0ze/golang-hft/src/internal"
	"github.com/d0ze/golang-hft/src/pkg/domain/entities"
	"github.com/d0ze/golang-hft/src/pkg/exchange"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// interface to implement a strategy in the application
type IStrategy interface {
	// the Open function takes as input the new candle and the current price trend
	// and if there is opportunity for a profit, returns a list of orders to be placed
	Open(trend entities.ITrend, candle entities.Candle, balance *entities.Balance, positions []*entities.Position) *entities.Order
	// the Close function takes as input the new candle, the current price trend and every
	// open position, and returns a list of orders to close the current positions with a profit
	// if possible
	Close(trend entities.ITrend, candle entities.Candle, positions []*entities.Position) *entities.Order
}

func CheckCost(order *entities.Order) bool {
	return order.GetMarketCost().GreaterThanOrEqual(entities.Markets.GetMinCost(order.Market))
}

func CheckVolume(order *entities.Order) bool {
	return order.InitialVolume.GreaterThanOrEqual(entities.Markets.GetOrderMin(order.Market))
}

func buildOpenOrder(
	market internal.Market,
	side internal.OrderSide,
	volume decimal.Decimal,
	price decimal.Decimal) *entities.Order {
	return &entities.Order{
		Id:            uuid.New().String(),
		Type:          internal.MARGIN,
		PriceType:     internal.MARKET,
		InitialVolume: volume,
		Side:          side,
		Status:        internal.CREATED,
		Market:        market,
		MarketPrice:   price,
		CreatedAt:     time.Now(),
		ReduceOnly:    false,
		Leverage:      exchange.KrakenCli.GetLeverage(market).BigInt().Int64(),
	}
}
func buildClosingOrder(position *entities.Position) *entities.Order {
	var side internal.OrderSide
	if position.Side == internal.BUY {
		side = internal.SELL
	} else {
		side = internal.BUY
	}
	return &entities.Order{
		Id:            uuid.New().String(),
		Type:          internal.SPOT,
		PriceType:     internal.MARKET,
		InitialVolume: position.Size,
		Side:          side,
		Status:        internal.CREATED,
		Market:        position.Market,
		CreatedAt:     time.Now(),
		ReduceOnly:    true,
		Leverage:      exchange.KrakenCli.GetLeverage(position.Market).BigInt().Int64(),
	}
}
