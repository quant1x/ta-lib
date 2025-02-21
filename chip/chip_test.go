package chip

import (
	"fmt"
	"gitee.com/quant1x/engine/datasource/base"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gotdx/securities"
	"testing"
	// 使用decimal库处理精确计算
	"github.com/shopspring/decimal"
)

func TestChisDecimal(t *testing.T) {
	low := 0.853
	price := decimal.NewFromFloat(low)
	price.Round(2).Float64()
	fmt.Println(price)

	v1, _ := decimal.NewFromFloat(9.824).Round(2).Float64()
	v2, _ := decimal.NewFromFloat(9.826).Round(2).Float64()
	v3, _ := decimal.NewFromFloat(9.8251).Round(2).Float64()
	fmt.Println(v1, v2, v3)

	v4, _ := decimal.NewFromFloat(9.815).Round(2).Float64()
	v5, _ := decimal.NewFromFloat(9.825).Round(2).Float64()
	v6, _ := decimal.NewFromFloat(9.835).Round(2).Float64()
	v7, _ := decimal.NewFromFloat(9.845).Round(2).Float64()
	fmt.Println(v4, v5, v6, v7)

	v8, _ := decimal.NewFromFloat(3.3).Round(2).Float64()
	v9, _ := decimal.NewFromFloat(3.3000000000000003).Round(2).Float64()
	v10, _ := decimal.NewFromFloat(3).Round(2).Float64()
	fmt.Println(v8, v9, v10)

	v11, _ := decimal.NewFromFloat(129.975).Round(2).Float64()
	v12, _ := decimal.NewFromFloat(34423.125).Round(2).Float64()
	fmt.Println(v11, v12)
}

func TestChips(t *testing.T) {
	// 示例配置
	config := defaultConfig

	// 初始化计算器
	calculator := NewChipDistribution(config)

	code := "300251"
	//code = "301256"
	//code = "300543"
	//code = "603980"
	//code = "600126"
	//code = "301487"
	code = "300098"
	//code = "002281"
	//code = "300170"
	date := "2025-02-20"
	//date = "2024-09-18"
	//date = exchange.GetFrontTradeDay()
	securityCode := exchange.CorrectSecurityCode(code)
	securityName := securities.GetStockName(securityCode)
	tradeDate := exchange.GetCurrentDate(date)
	fmt.Printf("%s(%s), 截至%s收盘: \n", securityName, securityCode, tradeDate)
	// 加载数据
	klines := base.CheckoutKLines(securityCode, tradeDate)
	if err := calculator.LoadCSV(securityCode, tradeDate); err != nil {
		fmt.Println("加载数据失败:", err)
		return
	}

	// 计算筹码分布
	if err := calculator.Calculate(); err != nil {
		fmt.Println("计算失败:", err)
		return
	}

	// 查找主要筹码峰
	targetPrice := klines[len(klines)-1].Close
	upper, lower, err := calculator.FindMainPeaks(targetPrice)
	if err != nil {
		fmt.Println("查找峰值失败:", err)
		return
	}

	fmt.Printf("当前价格 %.2f 附近的主要筹码峰:\n", targetPrice)
	fmt.Printf("压力: 最接近=%.2f, 最大=%.2f, 成交量=%.2f股, 占比=%.2f%%\n", upper.Closest, upper.Extremum, upper.Volume, 100*upper.Proportion)
	fmt.Printf("支撑: 最接近=%.2f, 最大=%.2f, 成交量=%.2f股, 占比=%.2f%%\n", lower.Closest, lower.Extremum, lower.Volume, 100*lower.Proportion)
}
