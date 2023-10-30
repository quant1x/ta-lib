package indicators

import (
	"fmt"
	"gitee.com/quant1x/engine/datasets"
	"testing"
)

func TestBRAR(t *testing.T) {
	df := datasets.KLine("sz002528")
	fmt.Println(df)
	df1 := BRAR(df, 26)
	fmt.Println(df1)
}
