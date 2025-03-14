package indicators

import (
	"fmt"
	"gitee.com/quant1x/engine/factors"
	"testing"
)

func TestPlatform(t *testing.T) {
	code := "600703.sh"
	code = "603789.sh"
	code = "sz000506"
	code = "sh603367"
	//code = "sz002275"
	code = "sz002665"
	code = "sz002528"
	//code = "sz000892"
	//code = "sz000905"
	code = "sh600641"
	code = "sh688031"
	code = "sz000988"
	code = "sh600105"
	code = "sz002292"
	code = "sh600354"
	code = "sh605577"
	code = "sh688662"
	code = "sz300678"
	code = "sh605162"
	code = "sz002992"
	df := factors.KLine(code)
	fmt.Println(df)
	df1 := Platform(df)
	fmt.Println(df1)
	_ = df1.WriteCSV("t02.csv")
}
