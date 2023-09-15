package sample

import (
	"fmt"
	"gitee.com/quant1x/ta-lib/testfiles"
	"testing"
)

func TestConfidenceInterval(t *testing.T) {
	df := testfiles.LoadTestData()
	fmt.Println(df)
	df = ConfidenceInterval(df, 5)
	fmt.Println(df)
}
