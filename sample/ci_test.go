package sample

import (
	"fmt"
	"testing"

	"gitee.com/quant1x/engine/factors"
)

func TestConfidenceInterval(t *testing.T) {
	code := "688351.sh"
	df := factors.KLine(code)
	fmt.Println(df)
	df = ConfidenceInterval(df, 5)
	fmt.Println(df)
}
