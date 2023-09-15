package indicators

import (
	"fmt"
	"gitee.com/quant1x/ta-lib/testfiles"
	"testing"
)

func TestW(t *testing.T) {
	df := testfiles.LoadTestData()
	//df = df.SelectRows(stat.RangeFinite(0, -5))
	fmt.Println(df)
	W(df, true, true)
}
