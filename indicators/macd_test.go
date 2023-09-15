package indicators

import (
	"fmt"
	"gitee.com/quant1x/ta-lib/testfiles"
	"testing"
)

func TestMACD(t *testing.T) {
	df := testfiles.LoadTestData()
	fmt.Println(df)
	df1 := MACD(df, 5, 13, 3)
	fmt.Println(df1)
}
