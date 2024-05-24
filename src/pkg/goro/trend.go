package goro

import (
	"sync"
	"time"

	"github.com/d0ze/golang-hft/src/internal"
	"github.com/d0ze/golang-hft/src/pkg/domain/entities"
	"github.com/d0ze/golang-hft/src/pkg/exchange"
	"github.com/sirupsen/logrus"
)

// polls ohlc and updates the given trend instance
// with the new candles for the timeframe
func PollOHLC(pair internal.Market, trend entities.ITrend, timeframe int, wg *sync.WaitGroup) chan entities.Candle {
	logrus.Infof("[%s] polling ohlc data (interval %d)", pair, timeframe)
	candles := make(chan entities.Candle)
	wg.Add(1)
	go func() {
		defer wg.Done()
		pollNext(pair, trend, candles, timeframe)
		tick := time.NewTicker(time.Duration(timeframe*30) * time.Second)
		for range tick.C {
			pollNext(pair, trend, candles, timeframe)
		}
	}()
	return candles
}

func pollNext(pair internal.Market, trend entities.ITrend, output chan entities.Candle, timeframe int) {
	candles, err := exchange.KrakenCli.GetOHLC(pair, timeframe)
	if err != nil {
		logrus.Warnf("[%s] error retrieving ohlc : %v", pair, err)
	}
	if len(candles) > 0 {
		next := candles[len(candles)-1]
		// skip updating if the last candle is equal
		frameCandles := *trend.GetCandles(timeframe)
		if len(frameCandles) > 0 && frameCandles[len(frameCandles)-1].Timestamp == next.Timestamp {
			return
		}
		trend.Update(next, timeframe)
		output <- next
	}
}
