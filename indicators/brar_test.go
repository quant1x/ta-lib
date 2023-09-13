package indicators

import (
	"fmt"
	"gitee.com/quant1x/stock/features"
	"testing"
)

func TestBRAR(t *testing.T) {
	df := features.KLine("sz002528")
	fmt.Println(df)
	df1 := BRAR(df, 26)
	fmt.Println(df1)
}
