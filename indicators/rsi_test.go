package indicators

import (
	"fmt"
	"gitee.com/quant1x/ta-lib/testfiles"
	"testing"
)

func TestRSI(t *testing.T) {
	df := testfiles.LoadTestData()
	fmt.Println(df)
	df1 := RSI(df, 6, 12, 24)
	fmt.Println(df1)
}
