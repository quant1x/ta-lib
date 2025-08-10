package chip

import (
	"fmt"
	"testing"

	"gitee.com/quant1x/engine/datasource/base"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/gotdx/quotes"
	"gitee.com/quant1x/gotdx/securities"
	"gitee.com/quant1x/num"
	"github.com/shopspring/decimal" // 使用decimal库处理精确计算
)

func TestChisDecimal(t *testing.T) {
	fmt.Println("9÷2+1 =", 9/2+1)
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

func totalMinutes(data []quotes.MinuteTime, minPrice, maxPrice float64) float64 {
	vol := 0.00
	for _, v := range data {
		currentPrice := num.Decimal(float64(v.Price))
		if currentPrice >= minPrice && currentPrice <= maxPrice {
			vol += float64(v.Vol)
		}
	}
	return vol
}

func TestChips(t *testing.T) {
	// 示例配置
	config := defaultConfig

	// 初始化计算器
	cd := NewChipDistribution(config)

	code := "300251"
	code = "600580"
	//code = "002195"
	//code = "600633"
	//code = "300076"
	code = "600392"
	code = "600633"
	code = "000701"
	//code = "603999"
	//code = "600725"
	//code = "002730"
	//code = "301171"
	//code = "002632"
	//code = "300017"
	//code = "300699"
	//code = "001914"
	code = "002730"
	code = "600348"
	code = "601001"
	code = "002741"
	code = "600398"
	code = "600580"
	code = "002332"
	code = "600580"
	code = "000521"
	code = "000158"
	//code = "300046"
	date := "2025-08-01"
	securityCode := exchange.CorrectSecurityCode(code)
	securityName := securities.GetStockName(securityCode)
	tradeDate := exchange.GetCurrentDate(date)
	fmt.Printf("%s(%s), 截至%s收盘: \n", securityName, securityCode, tradeDate)
	// 加载数据
	klines := base.CheckoutKLines(securityCode, tradeDate)
	if err := cd.LoadCSV(securityCode, tradeDate); err != nil {
		fmt.Println("加载数据失败:", err)
		return
	}

	// 计算筹码分布
	if err := cd.Calculate(); err != nil {
		fmt.Println("计算失败:", err)
		return
	}

	// 查找主要筹码峰
	targetPrice := cd.LastClose
	upper, lower, err := cd.FindMainPeaks(targetPrice)
	if err != nil {
		fmt.Println("查找峰值失败:", err)
		return
	}

	fmt.Printf("当前价格 %.2f 附近的主要筹码峰:\n", targetPrice)
	fmt.Printf("压力(上): 最接近=%.2f, 最大=%.2f, 成交量=%.2f股, 占比=%.2f%%\n", upper.Closest, upper.Extremum, cd.RealVolume(upper.CurrentToPeakRatio), 100*upper.CurrentToPeakRatio)
	fmt.Printf("支撑(下): 最接近=%.2f, 最大=%.2f, 成交量=%.2f股, 占比=%.2f%%\n", lower.Closest, lower.Extremum, cd.RealVolume(lower.CurrentToPeakRatio), 100*lower.CurrentToPeakRatio)
	// 计算短线是否获得支撑
	n := len(klines)
	prevBar := klines[n-2]
	lastBar := klines[n-1]
	fmt.Printf("%+v\n", prevBar)
	fmt.Printf("%+v\n", lastBar)
	if prevBar.Close > lastBar.Close && lastBar.Low <= lower.Extremum && lastBar.Close > lower.Extremum && lower.CurrentToPeakVol > upper.CurrentToPeakVol {
		fmt.Println("\t=> 短线止跌")
	}
	if prevBar.Close < lower.Closest && lastBar.Close > lower.Extremum && lastBar.Close > lower.Closest && lower.CurrentToPeakVol > upper.CurrentToPeakVol {
		fmt.Println("\t=> 短线突破")
	}
	upperVolume := cd.RealVolume(upper.CurrentToPeakRatio)
	tradeDate = exchange.NextTradeDate(tradeDate)
	fmt.Printf("交易日期: %s\n", tradeDate)
	minutes := base.GetMinutes(securityCode, tradeDate)
	todayVol := totalMinutes(minutes, upper.Closest, upper.Closest*1.01)
	fmt.Printf("压力位: 预计抛压=%.2f, 实际成交=%.2f\n", upperVolume, todayVol*100)
	if todayVol*100 > upperVolume {
		fmt.Println("\t=>放量突破")
	}
	todayVol = totalMinutes(minutes, lower.Closest, lower.Closest*1.01)
	lowerVolume := cd.RealVolume(lower.CurrentToPeakRatio)
	fmt.Printf("支撑位: 预计支撑=%.2f, 实际成交=%.2f\n", lowerVolume, todayVol*100)
	if todayVol*100 > lowerVolume {
		fmt.Println("\t=>获得强支撑")
	}
}
