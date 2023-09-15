package testfiles

import (
	"fmt"
	"gitee.com/quant1x/pandas"
	"os"
)

func LoadTestData() pandas.DataFrame {
	fmt.Println(os.Getwd())
	filename := "../testfiles/sh600105.csv"
	df := pandas.ReadCSV(filename)
	return df
}
