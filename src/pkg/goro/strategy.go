package goro

import (
	"sync"

	"github.com/d0ze/golang-hft/src/internal"
	"github.com/d0ze/golang-hft/src/pkg/domain/entities"
	"github.com/d0ze/golang-hft/src/pkg/exchange"
	"github.com/d0ze/golang-hft/src/pkg/strategy"
	"github.com/sirupsen/logrus"
)

// goroutine which applies the strategy on each new candle and fires every
// order to be open into the returned channel
func Check(strategy strategy.IStrategy, market internal.Market, trend entities.ITrend, candles chan entities.Candle, wg *sync.WaitGroup) chan *entities.Order {
	result := make(chan *entities.Order)
	go func() {
		for candle := range candles {
			balance, err := exchange.KrakenCli.GetBalance()
			if err != nil {
				logrus.Warnf("[%s] error %v retrieving balance, skipping...", err, market)
				continue
			}
			positions, err := exchange.KrakenCli.GetOpenPositions(market)
			if err != nil {
				logrus.Warnf("[%s] error %v retrieving positions, skipping...", market, err)
				continue
			}
			open := strategy.Open(trend, candle, balance, positions)
			logrus.Infof("[%s] selected open order: %v", market, open)

			close := strategy.Close(trend, candle, positions)
			logrus.Infof("[%s] selected closing order: %v", market, close)
			if close != nil {
				result <- close
			}
			if open != nil {
				result <- open
			}
		}
	}()
	return result
}
