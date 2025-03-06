package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"strconv"
)

// 股价数据点
type StockPrice struct {
	Date  string
	Close float64
}

// 抛物线拟合结果
type ParabolaFit struct {
	A         float64   // 二次项系数
	B         float64   // 一次项系数
	C         float64   // 常数项
	RSquared  float64   // R²拟合优度
	Projected []float64 // 预测价格序列
}

// 二次多项式回归
func QuadraticRegression(prices []float64) ParabolaFit {
	n := float64(len(prices))
	var sumT, sumT2, sumT3, sumT4, sumY, sumTY, sumT2Y float64

	for t := 0; t < len(prices); t++ {
		ft := float64(t)
		y := prices[t]

		sumT += ft
		sumT2 += ft * ft
		sumT3 += ft * ft * ft
		sumT4 += ft * ft * ft * ft
		sumY += y
		sumTY += ft * y
		sumT2Y += ft * ft * y
	}

	// 构建正规方程矩阵
	S := [3][3]float64{
		{n, sumT, sumT2},
		{sumT, sumT2, sumT3},
		{sumT2, sumT3, sumT4},
	}
	V := [3]float64{sumY, sumTY, sumT2Y}

	// 高斯消元法求解系数
	det := S[0][0]*(S[1][1]*S[2][2]-S[1][2]*S[2][1]) -
		S[0][1]*(S[1][0]*S[2][2]-S[1][2]*S[2][0]) +
		S[0][2]*(S[1][0]*S[2][1]-S[1][1]*S[2][0])

	a := (V[0]*(S[1][1]*S[2][2]-S[1][2]*S[2][1]) -
		V[1]*(S[0][1]*S[2][2]-S[0][2]*S[2][1]) +
		V[2]*(S[0][1]*S[1][2]-S[0][2]*S[1][1])) / det

	b := (S[0][0]*(V[1]*S[2][2]-V[2]*S[2][1]) -
		S[0][1]*(V[0]*S[2][2]-V[2]*S[2][0]) +
		S[0][2]*(V[0]*S[2][1]-V[1]*S[2][0])) / det

	c := (S[0][0]*(S[1][1]*V[2]-S[1][2]*V[1]) -
		S[0][1]*(S[1][0]*V[2]-S[1][2]*V[0]) +
		S[0][2]*(S[1][0]*V[1]-S[1][1]*V[0])) / det

	// 计算R²
	var ssTotal, ssResidual float64
	meanY := sumY / n
	for t := 0; t < len(prices); t++ {
		yHat := a*float64(t)*float64(t) + b*float64(t) + c
		ssTotal += math.Pow(prices[t]-meanY, 2)
		ssResidual += math.Pow(prices[t]-yHat, 2)
	}
	rSquared := 1 - ssResidual/ssTotal

	// 生成预测值
	projected := make([]float64, len(prices))
	for t := range projected {
		projected[t] = a*float64(t)*float64(t) + b*float64(t) + c
	}

	return ParabolaFit{A: a, B: b, C: c, RSquared: rSquared, Projected: projected}
}

func main() {
	// 示例数据（可替换为真实CSV读取）
	historicalData := []StockPrice{
		{"2023-01-03", 100.0},
		{"2023-01-04", 102.5},
		{"2023-01-05", 105.8},
		{"2023-01-06", 110.2},
		{"2023-01-07", 115.3},
		{"2023-01-10", 121.1},
		{"2023-01-11", 127.6},
	}

	// 提取收盘价序列
	prices := make([]float64, len(historicalData))
	for i := range historicalData {
		prices[i] = historicalData[i].Close
	}

	// 进行抛物线拟合
	fit := QuadraticRegression(prices)

	// 输出结果
	fmt.Printf("拟合方程: y = %.4ft² + %.4ft + %.4f\n", fit.A, fit.B, fit.C)
	fmt.Printf("拟合优度 R² = %.4f\n", fit.RSquared)
	fmt.Printf("预测顶点时间: Day %.1f\n", -fit.B/(2*fit.A))
	fmt.Printf("预测峰值价格: %.2f\n", fit.C-(fit.B*fit.B)/(4*fit.A))

	// 写入CSV
	file, _ := os.Create("price_fit.csv")
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Write([]string{"Date", "Actual", "Projected"})

	for i, p := range historicalData {
		record := []string{
			p.Date,
			strconv.FormatFloat(p.Close, 'f', 2, 64),
			strconv.FormatFloat(fit.Projected[i], 'f', 2, 64),
		}
		writer.Write(record)
	}
}
