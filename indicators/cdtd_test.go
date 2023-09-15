package indicators

import (
	"fmt"
	"gitee.com/quant1x/ta-lib/testfiles"
	"testing"
)

func TestCDTD(t *testing.T) {
	df := testfiles.LoadTestData()
	fmt.Println(df)
	df1 := CDTD(df)
	fmt.Println(df1)
}
