package linear

import (
	"fmt"
	"gitee.com/quant1x/engine/datasource/base"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/pandas"
	"testing"
)

func TestTrendLine(t *testing.T) {
	code := "600839"
	code = exchange.CorrectSecurityCode(code)
	date := exchange.GetCurrentlyDay()
	rawData := base.CheckoutKLines(code, date)
	df := pandas.LoadStructs(rawData)
	df = TrendLine(df)
	fmt.Println(df)
}

func TestCrossTrend(t *testing.T) {
	code := "600839"
	code = exchange.CorrectSecurityCode(code)
	date := exchange.GetCurrentlyDay()
	rawData := base.CheckoutKLines(code, date)
	df := pandas.LoadStructs(rawData)
	df = CrossTrend(df)
	fmt.Println(df)
}
