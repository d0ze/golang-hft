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
	GetTwap(timeframe int) *decimal.Decimal
	// return the candle at the given position
	GetCandle(position int, timeframe int) Candle
	// return all the latest candles
	GetCandles(timeframe int) *[]Candle
	// returrns the trend market
	GetMarket() internal.Market
	GetSMA(period int, timeframe int) *decimal.Decimal
	GetRSI(period int, timeframe int) *decimal.Decimal
	GetBB(period int, stdDev float64, timeframe int) (*decimal.Decimal, *decimal.Decimal, *decimal.Decimal)
	GetMACD(fastPeriod int, slowPeriod int, signalPeriod int, timeframe int) (*decimal.Decimal, *decimal.Decimal)
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
	candles := t.GetCandles(timeframe)
	if len(*candles) >= internal.Config.OHLCSize {
		*candles = (*candles)[1:]
	}
	*candles = append(*candles, new)
}

func (t *trend) GetTwap(timeframe int) *decimal.Decimal {
	candles := t.GetCandles(timeframe)
	if len(*candles) == 0 {
		return nil
	}
	precision := Markets.GetDecimals(t.market)
	var sumWeightedPrice decimal.Decimal
	var totalTimeWeight int

	for _, candle := range *candles {
		duration := time.Since(candle.Timestamp).Minutes()
		priceWeight := candle.Close.Mul(decimal.NewFromFloat(duration))
		sumWeightedPrice = sumWeightedPrice.Add(priceWeight)
		totalTimeWeight += int(duration)
	}

	if totalTimeWeight == 0 {
		return nil
	}

	twap := utils.MarketPrecision(sumWeightedPrice.Div(decimal.NewFromInt(int64(totalTimeWeight))), precision)
	return &twap
}

func (t *trend) GetSMA(period int, timeframe int) *decimal.Decimal {
	candles := t.GetCandles(timeframe)
	if len(*candles) < period {
		return nil
	}

	sum := decimal.Zero
	for _, candle := range (*candles)[len(*candles)-period:] {
		sum = sum.Add(candle.Close)
	}
	sma := sum.Div(decimal.NewFromInt(int64(period)))
	r := utils.MarketPrecision(sma, Markets.GetDecimals(t.market))
	return &r
}

func (t *trend) GetRSI(period int, timeframe int) *decimal.Decimal {
	candles := t.GetCandles(timeframe)
	if len(*candles) < period+1 {
		return nil
	}

	gains := decimal.Zero
	losses := decimal.Zero

	for i := 1; i <= period; i++ {
		change := (*candles)[len(*candles)-i].Close.Sub((*candles)[len(*candles)-i-1].Close)
		if change.GreaterThan(decimal.Zero) {
			gains = gains.Add(change)
		} else {
			losses = losses.Add(change.Abs())
		}
	}

	averageGain := gains.Div(decimal.NewFromInt(int64(period)))
	averageLoss := losses.Div(decimal.NewFromInt(int64(period)))

	if averageLoss.IsZero() {
		averageLoss = decimal.NewFromFloat(1.0)
	}

	rs := averageGain.Div(averageLoss)
	rsi := decimal.NewFromInt(100).Sub(decimal.NewFromInt(100).Div(decimal.NewFromInt(1).Add(rs)))

	return &rsi
}

func (t *trend) GetBB(period int, stdDev float64, timeframe int) (*decimal.Decimal, *decimal.Decimal, *decimal.Decimal) {
	sma := t.GetSMA(period, timeframe)
	if sma == nil {
		return nil, nil, nil
	}

	candles := t.GetCandles(timeframe)
	var sumSqDiff decimal.Decimal
	for _, candle := range *candles {
		diff := candle.Close.Sub(*sma)
		sumSqDiff = sumSqDiff.Add(diff.Mul(diff))
	}
	variance := sumSqDiff.Div(decimal.NewFromInt(int64(period)))
	stdDeviation := decimal.NewFromFloat(math.Sqrt(variance.InexactFloat64()))

	upperBand := sma.Add(stdDeviation.Mul(decimal.NewFromFloat(stdDev)))
	lowerBand := sma.Sub(stdDeviation.Mul(decimal.NewFromFloat(stdDev)))
	precision := Markets.GetDecimals(t.market)
	r1, r2, r3 := utils.MarketPrecision(upperBand, precision), utils.MarketPrecision(lowerBand, precision), utils.MarketPrecision(*sma, precision)
	return &r1, &r2, &r3
}

func (t *trend) GetMACD(fastPeriod int, slowPeriod int, signalPeriod int, timeframe int) (*decimal.Decimal, *decimal.Decimal) {
	candles := *t.GetCandles(timeframe)
	if len(candles) < slowPeriod {
		return nil, nil
	}
	precision := Markets.GetDecimals(t.market)

	emaFast := calculateEMA(candles, fastPeriod)
	emaSlow := calculateEMA(candles, slowPeriod)
	macd := emaFast.Sub(emaSlow)
	macdSignal := calculateSignalLine(macd, signalPeriod)

	r1, r2 := utils.MarketPrecision(macd, precision), utils.MarketPrecision(macdSignal, precision)
	return &r1, &r2
}

func calculateEMA(candles []Candle, period int) decimal.Decimal {
	k := decimal.NewFromFloat(2.0 / float64(period+1))
	ema := candles[0].Close
	for i := 1; i < len(candles); i++ {
		ema = candles[i].Close.Mul(k).Add(ema.Mul(decimal.NewFromInt(1).Sub(k)))
	}
	return ema
}

func calculateSignalLine(macd decimal.Decimal, period int) decimal.Decimal {
	k := decimal.NewFromFloat(2.0 / float64(period+1))
	signal := macd
	for i := 1; i < period; i++ {
		signal = macd.Mul(k).Add(signal.Mul(decimal.NewFromInt(1).Sub(k)))
	}
	return signal
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
