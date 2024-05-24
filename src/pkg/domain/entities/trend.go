package entities

import (
	"math"
	"time"

	"github.com/d0ze/golang-hft/src/internal"
	"github.com/d0ze/golang-hft/src/pkg/utils"
	"github.com/shopspring/decimal"
)

type Timeframe int

const (
	TIMEFRAME_1M  Timeframe = 1
	TIMEFRAME_5M  Timeframe = 5
	TIMEFRAME_15M Timeframe = 15
	TIMEFRAME_1H  Timeframe = 60
)

type ITrend interface {
	// adds a new candle price to the current trend
	// if already 12 candles are present, we also remove the
	// first price and shift the slice
	Update(new Candle, timeframe int)
	// returns the time-weighted avg price of the last
	// 12 candles
	GetTwap(timeframe int) decimal.Decimal
	// return the candle at the given position
	GetCandle(position int, timeframe int) Candle
	// return all the latest candles
	GetCandles(timeframe int) *[]Candle
	// returrns the trend market
	GetMarket() internal.Market
	GetSMA(period int, timeframe int) decimal.Decimal
	GetRSI(period int, timeframe int) decimal.Decimal
	GetBB(period int, stdDev float64, timeframe int) (decimal.Decimal, decimal.Decimal, decimal.Decimal)
}

type trend struct {
	timeframes map[Timeframe]*[]Candle
	market     internal.Market
}

func InitTrend(market internal.Market) ITrend {
	return &trend{market: market, timeframes: map[Timeframe]*[]Candle{
		TIMEFRAME_1M:  {},
		TIMEFRAME_5M:  {},
		TIMEFRAME_15M: {},
		TIMEFRAME_1H:  {},
	}}
}

func (t *trend) Update(new Candle, timeframe int) {
	candles := *t.timeframes[Timeframe(timeframe)]
	if len(candles) == internal.Config.OHLCSize {
		candles = candles[1:]
	}
	candles = append(candles, new)
	t.timeframes[Timeframe(timeframe)] = &candles
}

// returns the time weighted average price of the trend
func (t *trend) GetTwap(timeframe int) decimal.Decimal {
	candles := *t.GetCandles(timeframe)
	if len(candles) == 0 {
		return decimal.Zero
	}
	precision := Markets.GetDecimals(t.market)
	var sumWeightedPrice decimal.Decimal
	var totalTimeWeight int // Total time weight (sum of time intervals)

	// Iterate over the last n candles (or available candles if less than n)
	for _, candle := range candles {
		duration := time.Since(candle.Timestamp).Minutes() // Calculate duration in minutes
		priceWeight := candle.Close.Mul(decimal.NewFromFloat(duration))
		sumWeightedPrice = sumWeightedPrice.Add(priceWeight)
		totalTimeWeight += int(duration)
	}

	if totalTimeWeight == 0 {
		return decimal.Zero
	}

	return utils.MarketPrecision(sumWeightedPrice.Div(decimal.NewFromInt(int64(totalTimeWeight))), precision)
}

// SMA calculates the Simple Moving Average over a specified period
func (t *trend) GetSMA(period int, timeframe int) decimal.Decimal {
	res := decimal.Zero
	precision := Markets.GetDecimals(t.market)
	candles := *t.GetCandles(timeframe)
	for _, candle := range candles {
		res = res.Add(candle.Close)
	}
	return utils.MarketPrecision(res.Div(decimal.NewFromInt(int64(len(candles)))), precision)
}

// RSI calculates the Relative Strength Index over a specified period
func (t *trend) GetRSI(period int, timeframe int) decimal.Decimal {
	candles := *t.GetCandles(timeframe)
	if len(candles) < period {
		return decimal.Zero
	}
	gain := decimal.Zero
	loss := decimal.Zero

	for i := 1; i < period+1; i++ {
		change := candles[i].Close.Sub(candles[i-1].Close)
		if change.GreaterThan(decimal.Zero) {
			gain = gain.Add(change)
		} else {
			loss = loss.Add(change.Abs())
		}
	}

	averageGain := gain.Div(decimal.NewFromInt(int64(period)))
	averageLoss := loss.Div(decimal.NewFromInt(int64(period)))

	if averageLoss.IsZero() {
		// Prevent division by zero
		return decimal.NewFromInt(100)
	}

	rs := averageGain.Div(averageLoss)
	rsi := decimal.NewFromInt(100).Sub(decimal.NewFromInt(100).Div(decimal.NewFromInt(1).Add(rs)))

	return rsi
}

// BB calculates the Bollinger Bands (middle, upper, lower) over a specified period and standard deviation
func (t *trend) GetBB(period int, stdDev float64, timeframe int) (decimal.Decimal, decimal.Decimal, decimal.Decimal) {
	sma := t.GetSMA(period, timeframe)
	if sma.IsZero() {
		return decimal.Zero, decimal.Zero, decimal.Zero
	}
	precision := Markets.GetDecimals(t.market)

	var sumSqDiff decimal.Decimal
	for _, candle := range *t.GetCandles(timeframe) {
		diff := candle.Close.Sub(sma)
		sumSqDiff = sumSqDiff.Add(diff.Mul(diff))
	}
	stdDeviationF, _ := decimal.NewFromFloat(stdDev).Mul(decimal.NewFromInt(int64(period))).Float64()
	stdDeviation := decimal.NewFromFloat(math.Sqrt(stdDeviationF))
	upperBand := sma.Add(stdDeviation)
	lowerBand := sma.Sub(stdDeviation)
	return utils.MarketPrecision(sma, precision), utils.MarketPrecision(upperBand, precision), utils.MarketPrecision(lowerBand, precision)
}

func (t *trend) GetCandle(position int, timeframe int) Candle {
	if len(*t.GetCandles(timeframe)) >= position {
		candles := *t.GetCandles(timeframe)
		return candles[position]
	}
	return Candle{}
}

func (t *trend) GetCandles(timeframe int) *[]Candle {
	return t.timeframes[Timeframe(timeframe)]
}

func (t *trend) GetMarket() internal.Market {
	return t.market
}
