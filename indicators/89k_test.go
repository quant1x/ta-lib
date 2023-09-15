package indicators

import (
	"fmt"
	"gitee.com/quant1x/ta-lib/testfiles"
	"testing"
)

func TestF89K(t *testing.T) {
	df := testfiles.LoadTestData()
	fmt.Println(df)
	df1 := F89K(df, 89)
	fmt.Println(df1)
}
