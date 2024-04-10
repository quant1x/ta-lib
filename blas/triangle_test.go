package blas

import (
	"fmt"
	"gitee.com/quant1x/engine/datasource/base"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gotdx/securities"
	"testing"
)

func TestChartKLine_Triangle(t *testing.T) {
	// 楔形模型测试
	requiredKLines := 89
	requiredKLines = 34
	//requiredKLines = 50
	//requiredKLines = 250
	code := "sh000001"
	date := "2024-04-10"
	//date = cache.DefaultCanReadDate()
	list := base.CheckoutKLines(code, date)
	if len(list) >= requiredKLines {
		list = list[len(list)-requiredKLines:]
	}
	sample := LoadKLineSample(list)
	securityCode := exchange.CorrectSecurityCode(code)
	waves := PeaksAndValleys(sample, securityCode)
	fmt.Println(waves)

	chartName := securities.GetStockName(code) + "(" + securityCode + ")日线图 - " + date
	graph := sample.Chart(chartName)
	//var pattern Pattern
	pattern := MatchTriangle(waves)
	if pattern != nil {
		fmt.Printf("triangle=%+v\n", pattern)
		series := pattern.Export(waves.Data)
		graph.AddSeries(series...)
	}
	name := "triangle-kline-" + securityCode
	_ = graph.Output(name)
}
