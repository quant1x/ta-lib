package linear

import (
	"fmt"
	"gitee.com/quant1x/engine/factors"
	"testing"
)

func TestTrendLine(t *testing.T) {
	code := "sh000905"
	code = "sz002528"
	//code = "sz002322"
	df := factors.KLine(code)
	df = TrendLine(df)
	fmt.Println(df)
}

func TestCrossTrend(t *testing.T) {
	code := "sh000905"
	code = "sz002528"
	//code = "sz002322"
	code = "sh600018"
	code = "sh603130"
	code = "sz002209"
	code = "sh600178"
	df := factors.KLine(code)
	df = CrossTrend(df)
	fmt.Println(df)
}
