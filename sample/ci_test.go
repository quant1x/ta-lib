package sample

import (
	"fmt"
	"gitee.com/quant1x/engine/factors"
	"testing"
)

func TestConfidenceInterval(t *testing.T) {
	code := "688351.sh"
	df := factors.KLine(code)
	fmt.Println(df)
	df = ConfidenceInterval(df, 5)
	fmt.Println(df)
}
