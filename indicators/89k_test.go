package indicators

import (
	"fmt"
	"gitee.com/quant1x/engine/factors"
	"testing"
)

func TestF89K(t *testing.T) {
	df := factors.KLine("sh600496")
	fmt.Println(df)
	df1 := F89K(df, 89)
	fmt.Println(df1)
}
