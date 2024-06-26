package indicators

import (
	"fmt"
	"gitee.com/quant1x/engine/datasource/base"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/pandas"
	"testing"
)

func TestW(t *testing.T) {
	code := "sh000905"
	code = "sz002528"
	//code = "sz000151"
	//code = "sz002564"
	//code = "sz002209"
	//code = "sz002951"
	//code = "sh000001"
	code = "sh600703"
	code = "sh688358"
	code = "002564.sz"
	code = "sh000001"
	code = "sz002728"
	code = "sh603528"
	code = "000888"
	code = "sh000001"
	//code = "sh600178"
	//code = "603066"
	date := "2024-06-25"
	code = exchange.CorrectSecurityCode(code)
	date = exchange.FixTradeDate(date)
	rows := base.CheckoutKLines(code, date)
	df := pandas.LoadStructs(rows)
	//df = df.SelectRows(stat.RangeFinite(0, -5))
	fmt.Println(df)
	fmt.Printf("   证券代码: %s\n", code)
	W(df, true, true)
}
