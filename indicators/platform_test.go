package indicators

import (
	"fmt"
	"gitee.com/quant1x/ta-lib/testfiles"
	"testing"
)

func TestPlatform(t *testing.T) {
	df := testfiles.LoadTestData()
	fmt.Println(df)
	df1 := Platform(df)
	fmt.Println(df1)
	_ = df1.WriteCSV("t02.csv")
}
