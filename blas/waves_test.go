package blas

import (
	"encoding/json"
	"fmt"
	"testing"

	"gitee.com/quant1x/engine/datasource/base"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gox/api"
)

func TestWaves_basic(t *testing.T) {
	requiredKLines := 0
	//requiredKLines = 250
	code := "sh000001"
	code = "sz000158"
	date := "2024-04-01"
	date = "2025-06-23"
	//date = cache.DefaultCanReadDate()
	list := base.CheckoutKLines(code, date)
	if requiredKLines > 0 && len(list) >= requiredKLines {
		list = list[len(list)-requiredKLines:]
	}
	sample := LoadKLineSample(list)
	//rows := sample.Len()
	securityCode := exchange.CorrectSecurityCode(code)
	waves := PeaksAndValleys(sample, securityCode)
	fmt.Println(waves)
	data, _ := json.Marshal(waves)
	text := api.Bytes2String(data)
	fmt.Println(text)
}
