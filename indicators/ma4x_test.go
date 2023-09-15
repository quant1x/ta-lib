package indicators

import (
	"fmt"
	"gitee.com/quant1x/ta-lib/testfiles"
	"testing"
)

func TestMA4X(t *testing.T) {
	df := testfiles.LoadTestData()
	fmt.Println(df)
	df1 := MA4X(df, 5)
	fmt.Println(df1)
}
