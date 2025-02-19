package chip

import (
	"fmt"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gotdx/securities"
	"testing"
)

func TestChips(t *testing.T) {
	// 示例配置
	config := defaultConfig

	// 初始化计算器
	calculator := NewChipDistribution(config)

	code := "300251"
	//code = "300271"
	//code = "300255"
	//code = "002539"
	code = "301238"
	code = "300173"
	code = "300759"
	date := "2025-02-18"
	securityCode := exchange.CorrectSecurityCode(code)
	securityName := securities.GetStockName(securityCode)
	tradeDate := exchange.GetCurrentDate(date)
	fmt.Printf("%s(%s), 截至%s收盘: \n", securityName, securityCode, tradeDate)
	// 加载数据
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
	targetPrice := 27.80
	//targetPrice = 10.01
	//targetPrice = 19.02
	targetPrice = 20.05
	targetPrice = 7.21
	targetPrice = 26.73
	upper, lower, err := calculator.FindMainPeaks(targetPrice)
	if err != nil {
		fmt.Println("查找峰值失败:", err)
		return
	}

	fmt.Printf("当前价格 %.2f 附近的主要筹码峰:\n", targetPrice)
	fmt.Printf("压力: 最接近=%.2f, 最大=%.2f\n", upper.Closest, upper.Extremum)
	fmt.Printf("支撑: 最接近=%.2f, 最大=%.2f\n", lower.Closest, lower.Extremum)
}
