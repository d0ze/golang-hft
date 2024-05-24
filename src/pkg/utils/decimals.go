package utils

import (
	"github.com/shopspring/decimal"
)

func Percentage(number decimal.Decimal, percent int) decimal.Decimal {
	return number.Mul(decimal.NewFromInt(int64(percent))).Div(decimal.NewFromInt(100))
}

func MarketPrecision(number decimal.Decimal, decimals int) decimal.Decimal {
	return decimal.RequireFromString(number.StringFixed(int32(decimals)))
}
