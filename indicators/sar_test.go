package indicators

import (
	"fmt"
	"gitee.com/quant1x/engine/datasource/base"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/pandas"
	"testing"
)

func TestSar_basic(t *testing.T) {
	code := "600171"
	date := "2024-06-13"
	code = exchange.CorrectSecurityCode(code)
	date = exchange.FixTradeDate(date)
	list := base.CheckoutKLines(code, date)
	rows := len(list)
	high := make([]float64, rows)
	low := make([]float64, rows)
	firstIsBull := false
	for i, v := range list {
		if i == 0 {
			firstIsBull = v.Close > v.Open
		}
		high[i] = v.High
		low[i] = v.Low
	}
	data := SAR(firstIsBull, high, low, 0.02, 0.20)
	df := pandas.LoadStructs(data)
	fmt.Println(df)
}
