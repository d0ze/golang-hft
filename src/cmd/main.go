package main

import (
	"strconv"
	"strings"
	"sync"

	"github.com/d0ze/golang-hft/src/internal"
	"github.com/d0ze/golang-hft/src/pkg/domain/entities"
	"github.com/d0ze/golang-hft/src/pkg/exchange"
	"github.com/d0ze/golang-hft/src/pkg/goro"
	"github.com/d0ze/golang-hft/src/pkg/strategy"
	"github.com/sirupsen/logrus"
)

func main() {
	// init internals

	internal.InitConfig()
	internal.InitLogging()
	strategies := map[string]strategy.IStrategy{
		"simple": strategy.NewSimpleStrategy(),
	}

	// init broker
	exchange.InitClient()
	markets := internal.Config.GetMarkets()
	data, err := exchange.KrakenCli.GetMarketsData(markets)
	if err != nil {
		logrus.Warnf("[MAIN] couldnt retrieve market data (reason: %v), exiting", err)
		return
	}

	entities.Markets = data
	logrus.Infof("[MAIN] retrieved markets data: %s", entities.Markets.String())

	logrus.Infof("[MAIN] selected strategy %s", internal.Config.Strategy)
	stategy := strategies[internal.Config.Strategy]

	var wg sync.WaitGroup

	for _, market := range markets {
		var ticks chan entities.Candle
		trend := entities.InitTrend(market)
		timeframes := strings.Split(internal.Config.OHLCIntervals, "-")
		logrus.Infof("[MAIN] selected timeframes %v", timeframes)
		for _, tf := range timeframes {
			// get ohlc for each timeframe
			timeframe, _ := strconv.ParseInt(tf, 10, 16)
			logrus.Infof("[MAIN] retrieving candles for timeframe %dm", timeframe)
			prev, err := exchange.KrakenCli.GetOHLC(market, int(timeframe))
			logrus.Infof("[MAIN] received  %d candles", len(prev))
			if err != nil {
				logrus.Fatalf("[MAIN] error %v retrieving latest %dm candles", timeframe, err)
			}
			for _, candle := range prev[len(prev)-60 : len(prev)-1] {
				logrus.Infof("[MAIN] loading candle %s", candle.String())
				trend.Update(candle, int(timeframe))
			}
			if int(timeframe) == internal.Config.StrategyIntervalCheck {
				// defines which candle timeframe will tick the strategy check
				ticks = goro.PollOHLC(market, trend, int(timeframe), &wg)
			}
		}
		orders := goro.Check(stategy, market, trend, ticks, &wg)
		goro.HandleOrders(orders)
	}
	wg.Wait()
	logrus.Infof("[MAIN] program exiting...")
}
