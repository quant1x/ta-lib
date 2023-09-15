package indicators

import (
	"fmt"
	"gitee.com/quant1x/ta-lib/testfiles"
	"testing"
)

func TestBRAR(t *testing.T) {
	df := testfiles.LoadTestData()
	fmt.Println(df)
	df1 := BRAR(df, 26)
	fmt.Println(df1)
}
