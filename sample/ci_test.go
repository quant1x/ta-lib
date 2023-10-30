package sample

import (
	"fmt"
	"gitee.com/quant1x/engine/datasets"
	"testing"
)

func TestConfidenceInterval(t *testing.T) {
	code := "688351.sh"
	df := datasets.KLine(code)
	fmt.Println(df)
	df = ConfidenceInterval(df, 5)
	fmt.Println(df)
}
