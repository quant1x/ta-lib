package indicators

import (
	"fmt"
	"testing"

	"gitee.com/quant1x/engine/datasource/base"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gotdx/securities"
	"gitee.com/quant1x/pandas"
)

func TestM(t *testing.T) {
	code := "sh000001"
	//code = "sh600178"
	//code = "603066"
	//code := "300781"
	date := "2024-06-25"
	code = exchange.CorrectSecurityCode(code)
	date = exchange.FixTradeDate(date)
	rows := base.CheckoutKLines(code, date)
	df := pandas.LoadStructs(rows)
	//df = df.SelectRows(stat.RangeFinite(0, -5))
	fmt.Println(df)
	fmt.Printf("   证券代码: %s, %s\n", code, securities.GetStockName(code))
	M(df, true, true)
}
