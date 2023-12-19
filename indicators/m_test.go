package indicators

import (
	"fmt"
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/gotdx/securities"
	"testing"
)

func TestM(t *testing.T) {
	code := "sh000001"
	df := factors.KLine(code)
	//df = df.SelectRows(stat.RangeFinite(0, -5))
	fmt.Println(df)
	fmt.Printf("   证券代码: %s, %s\n", code, securities.GetStockName(code))
	M(df, true, true)
}
