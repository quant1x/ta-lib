package blas

import (
	"fmt"
	"testing"

	"github.com/quant1x/engine/datasource/base"
	"github.com/quant1x/exchange"
	"github.com/quant1x/gotdx/securities"
)

func TestWedge_basic(t *testing.T) {
	// 楔形模型测试
	requiredKLines := 89
	requiredKLines = 34
	//requiredKLines = 50
	//requiredKLines = 250
	code := "sh000001"
	//code = "300824"
	//code = "300945"
	//code = "300107"
	//code = "002242"
	//code = "600855"
	//code = "300019"
	//code = "300107"
	//code = "600719"
	code = "sz000158"
	date := "2024-04-01"
	date = "2025-07-10"
	//date = "2024-03-29"
	//date = cache.DefaultCanReadDate()
	date = exchange.FixTradeDate(date)
	list := base.CheckoutKLines(code, date)
	if len(list) >= requiredKLines {
		list = list[len(list)-requiredKLines:]
	}
	sample := LoadKLineSample(list)
	securityCode := exchange.CorrectSecurityCode(code)
	waves := HighAndLow(sample, securityCode)
	fmt.Println(waves)

	chartName := securities.GetStockName(code) + "(" + securityCode + ")日线图 - " + date
	graph := sample.Chart(chartName)
	var pattern Pattern
	pattern = MatchWedge(waves)
	if pattern != nil {
		wedge := pattern.(*Wedge)
		dp, top, bottom, ok := wedge.Cross()
		fmt.Println(dp, top, bottom, ok)
		fmt.Printf("wedge=%+v\n", pattern)
		series := pattern.ExportSeries(sample)
		graph.AddSeries(series...)
		series = pattern.NeckSeries(sample)
		graph.AddSeries(series...)
	}
	name := "wedge-kline-" + securityCode
	_ = graph.Output(name)
}
