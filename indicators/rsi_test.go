package indicators

import (
	"fmt"
	"gitee.com/quant1x/stock/features"
	"testing"
)

func TestRSI(t *testing.T) {
	df := features.KLine("sz002528")
	fmt.Println(df)
	df1 := RSI(df, 6, 12, 24)
	fmt.Println(df1)
}
