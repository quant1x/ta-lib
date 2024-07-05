package indicators

import (
	"fmt"
	"gitee.com/quant1x/engine/datasource/base"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/pandas"
	"testing"
)

func TestMACD(t *testing.T) {
	code := "300781"
	code = "002766"
	date := "2024-07-02"
	code = exchange.CorrectSecurityCode(code)
	date = exchange.FixTradeDate(date)
	rows := base.CheckoutKLines(code, date)
	df := pandas.LoadStructs(rows)
	fmt.Println(df)
	df1 := MACD(df, 12, 26, 9)
	fmt.Println(df1)
}

func TestMinutesMacd(t *testing.T) {
	code := "sh000001"
	code = "sh510050"
	code = "sh600105"
	code = "880948"
	code = "603228"
	code = "002766"
	date := "2024-07-03"
	code = exchange.CorrectSecurityCode(code)
	date = exchange.FixTradeDate(date)
	v := base.GetMinutes(code, date)
	df := pandas.LoadStructs(v)
	prices := df.ColAsNDArray("Price")
	df1 := macd(prices, 12, 26, 9)
	fmt.Println(df1)
}
