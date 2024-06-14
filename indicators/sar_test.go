package indicators

import (
	"fmt"
	"gitee.com/quant1x/engine/datasource/base"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/pandas"
	"testing"
)

func TestSar_basic(t *testing.T) {
	code := "300046"
	date := "2024-05-27"
	code = exchange.CorrectSecurityCode(code)
	date = exchange.FixTradeDate(date)
	list := base.CheckoutKLines(code, date)
	rows := len(list)
	high := make([]float64, rows)
	low := make([]float64, rows)
	for i, v := range list {
		high[i] = v.High
		low[i] = v.Low
	}
	data := SAR(high, low)
	df := pandas.LoadStructs(data)
	fmt.Println(df)
	last := data[rows-1]
	// 增量计算SAR 2024-06-14
	latest := last.Incr(18.77, 17.88)
	fmt.Printf("%+v\n", latest)
}
