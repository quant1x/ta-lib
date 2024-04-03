package blas

import (
	"fmt"
	"gitee.com/quant1x/engine/datasource/base"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gotdx/securities"
	"testing"
)

func TestWedge_basic(t *testing.T) {
	// 楔形模型测试
	requiredKLines := 89
	//requiredKLines = 34
	//requiredKLines = 50
	//requiredKLines = 250
	code := "sh000001"
	//code = "sz300629"
	//code = "000917"
	//code = "600843"
	//code = "sz000751"
	code = "600603"
	date := "2024-04-02"
	//date = "2024-03-29"
	//date = cache.DefaultCanReadDate()
	list := base.CheckoutKLines(code, date)
	if len(list) >= requiredKLines {
		list = list[len(list)-requiredKLines:]
	}
	sample := LoadKLineSample(list)
	securityCode := exchange.CorrectSecurityCode(code)
	waves := NewWaves(sample, securityCode)
	fmt.Println(waves)

	chartName := securities.GetStockName(code) + "(" + securityCode + ")日线图 - " + date
	graph := sample.Chart(chartName)
	var pattern Pattern
	pattern = MatchWedge(waves)
	if pattern != nil {
		fmt.Printf("wedge=%+v\n", pattern)
		series := pattern.ExportSeries(sample)
		graph.AddSeries(series...)
	}
	name := "wedge-kline-" + securityCode
	_ = graph.Output(name)
}
