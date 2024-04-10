package blas

import (
	"fmt"
	"gitee.com/quant1x/engine/datasource/base"
	"gitee.com/quant1x/exchange"
	"testing"
)

func TestWaves_basic(t *testing.T) {
	requiredKLines := 89
	//requiredKLines = 250
	code := "sh000001"
	code = "sz300629"
	date := "2024-04-01"
	date = "2024-03-29"
	//date = cache.DefaultCanReadDate()
	list := base.CheckoutKLines(code, date)
	if len(list) >= requiredKLines {
		list = list[len(list)-requiredKLines:]
	}
	sample := LoadKLineSample(list)
	//rows := sample.Len()
	securityCode := exchange.CorrectSecurityCode(code)
	waves := PeaksAndValleys(sample, securityCode)
	fmt.Println(waves)
}
