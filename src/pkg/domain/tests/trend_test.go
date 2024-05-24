package tests

import (
	"testing"
	"time"

	"github.com/d0ze/golang-hft/src/internal"
	"github.com/d0ze/golang-hft/src/pkg/domain/entities"
	"github.com/shopspring/decimal"
)

func TestSMA(t *testing.T) {
	// Create a mock trend with predefined candles
	internal.InitConfig()
	internal.InitLogging()
	markets := entities.NewMarkets()
	markets.SetMetadata(
		internal.XBTEUR,
		8,
		internal.XBT,
		internal.EUR,
		decimal.RequireFromString("0.01"),
		decimal.RequireFromString("0.00000001"),
	)
	entities.Markets = markets
	trend := entities.InitTrend(internal.XBTEUR)
	candles := []entities.Candle{
		{Timestamp: time.Now().Add(-3 * time.Hour), Close: decimal.NewFromFloat(100)},
		{Timestamp: time.Now().Add(-2 * time.Hour), Close: decimal.NewFromFloat(110)},
		{Timestamp: time.Now().Add(-1 * time.Hour), Close: decimal.NewFromFloat(120)},
		{Timestamp: time.Now(), Close: decimal.NewFromFloat(130)},
	}
	for _, candle := range candles {
		trend.Update(candle, 60)
	}

	// Calculate SMA
	sma := trend.GetSMA(4, 60)

	// Expected value is calculated manually
	expectedSMA := decimal.NewFromFloat(115)

	// Compare expected value with calculated value
	if !expectedSMA.Equal(*sma) {
		t.Errorf("SMA calculation error. Expected: %s, Got: %s", expectedSMA.String(), sma.String())
	}
}

func TestRSI(t *testing.T) {
	// Create a mock trend with predefined candles
	internal.InitConfig()
	internal.InitLogging()
	markets := entities.NewMarkets()
	markets.SetMetadata(
		internal.XBTEUR,
		8,
		internal.XBT,
		internal.EUR,
		decimal.RequireFromString("0.01"),
		decimal.RequireFromString("0.00000001"),
	)
	entities.Markets = markets
	trend := entities.InitTrend(internal.XBTEUR)
	candles := []entities.Candle{
		{Close: decimal.NewFromFloat(100)},
		{Close: decimal.NewFromFloat(110)},
		{Close: decimal.NewFromFloat(130)},
		{Close: decimal.NewFromFloat(120)},
		{Close: decimal.NewFromFloat(105)},
	}
	for _, candle := range candles {
		trend.Update(candle, 60)
	}

	rsi := trend.GetRSI(4, 60)

	expectedRSI := decimal.NewFromFloat(54.54)

	if !expectedRSI.Equal(*rsi) {
		t.Errorf("RSI calculation error. Expected: %s, Got: %s", expectedRSI.String(), rsi.String())
	}
}

func TestBollingerBands(t *testing.T) {
	// Create a mock trend with predefined candles
	internal.InitConfig()
	internal.InitLogging()
	markets := entities.NewMarkets()
	markets.SetMetadata(
		internal.XBTEUR,
		8,
		internal.XBT,
		internal.EUR,
		decimal.RequireFromString("0.01"),
		decimal.RequireFromString("0.00000001"),
	)
	entities.Markets = markets
	trend := entities.InitTrend(internal.XBTEUR)
	candles := []entities.Candle{
		{Close: decimal.NewFromFloat(100)},
		{Close: decimal.NewFromFloat(110)},
		{Close: decimal.NewFromFloat(120)},
		{Close: decimal.NewFromFloat(130)},
		{Close: decimal.NewFromFloat(140)},
	}
	for _, candle := range candles {
		trend.Update(candle, 60)
	}

	// Calculate Bollinger Bands
	upper, lower, middle := trend.GetBB(5, 2, 60)

	// Expected values are calculated manually
	expectedUpper := decimal.NewFromFloat(148.28427125)
	expectedMiddle := decimal.NewFromFloat(120)
	expectedLower := decimal.NewFromFloat(91.71572875)

	// Compare expected values with calculated values
	if !expectedUpper.Equal(*upper) {
		t.Errorf("Upper Bollinger Band calculation error. Expected: %s, Got: %s", expectedUpper.String(), upper.String())
	}
	if !expectedMiddle.Equal(*middle) {
		t.Errorf("Middle Bollinger Band calculation error. Expected: %s, Got: %s", expectedMiddle.String(), middle.String())
	}
	if !expectedLower.Equal(*lower) {
		t.Errorf("Lower Bollinger Band calculation error. Expected: %s, Got: %s", expectedLower.String(), lower.String())
	}
}
