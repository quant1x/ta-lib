package indicators

import (
	"fmt"
	"gitee.com/quant1x/engine/datasets"
	"testing"
)

func TestKDJ(t *testing.T) {
	df := datasets.KLine("sz002528")
	fmt.Println(df)
	df1 := KDJ(df, 9, 3, 3)
	fmt.Println(df1)
}
