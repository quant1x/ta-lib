package indicators

import (
	"fmt"
	"gitee.com/quant1x/stock/features"
	"testing"
)

func TestMACD(t *testing.T) {
	df := features.KLine("sz002528")
	fmt.Println(df)
	df1 := MACD(df, 5, 13, 3)
	fmt.Println(df1)
}
