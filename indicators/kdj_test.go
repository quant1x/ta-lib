package indicators

import (
	"fmt"
	"gitee.com/quant1x/ta-lib/testfiles"
	"testing"
)

func TestKDJ(t *testing.T) {
	df := testfiles.LoadTestData()
	fmt.Println(df)
	df1 := KDJ(df, 9, 3, 3)
	fmt.Println(df1)
}
