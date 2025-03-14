package indicators

import (
	"fmt"
	"gitee.com/quant1x/engine/factors"
	"testing"
)

func TestMA4X(t *testing.T) {
	code := "000736.sz"
	df := factors.KLine(code)
	fmt.Println(df)
	df1 := MA4X(df, 5)
	fmt.Println(df1)
}
