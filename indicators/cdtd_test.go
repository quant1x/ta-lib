package indicators

import (
	"fmt"
	"gitee.com/quant1x/stock/features"
	"testing"
)

func TestCDTD(t *testing.T) {
	code := "002528.sz"
	code = "sh600025"
	df := features.KLine(code)
	fmt.Println(df)
	df1 := CDTD(df)
	fmt.Println(df1)
}
