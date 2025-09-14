package blas

import (
	"cmp"

	"github.com/quant1x/num"
	"github.com/quant1x/ta-lib/plot"
)

// TendencyDirection 交易方向
type TendencyDirection = uint8

const (
	TradingKeep              TendencyDirection = 0x00   // 持股
	TradingBuy               TendencyDirection = 1 << 6 // 买入,01xxxxxx
	TradingSell              TendencyDirection = 1 << 7 // 卖出,10xxxxxx
	TradingBottomDivergence  TendencyDirection = 0x01   // 底背离
	TradingTopDivergence     TendencyDirection = 0x02   // 顶背离
	TradingConsistentRise    TendencyDirection = 0x03   // 一致性上涨
	TradingConsistentDecline TendencyDirection = 0x04   // 一致性下跌
)

// WavesTendency 波段趋势
type WavesTendency = int

const (
	TendencyPrice          WavesTendency = 0 // 股价主导
	TendencyLinear         WavesTendency = 1 // 线性趋势主导
	TendencyPriceAndLinear WavesTendency = 2 // 股价和趋势并存
)

// OperationalSignal 操作信号
type OperationalSignal struct {
	Signal TendencyDirection // 操作信号
	Digits int               // 保留小数点几位
}

// Analysis 分析接口
type Analysis interface {
	// Normalize 归一化处理(Initialize)
	Normalize() error
	// Process 加工处理
	Process() error
	// At 获取数据, 返回即时值, 最高值, 最低值
	At(n int) (current, high, low float64)
	// Match 匹配
	Match()
}

// Pattern 模式, 形态
type Pattern interface {
	// Fit 拟合, 输出支撑线, 颈线, 压力线
	Fit() (neckLine, supportLine, pressureLine num.Line)
	// ExportSeries 输出图表
	ExportSeries(sample DataSample) []plot.Series
	// NeckSeries 输出颈线
	NeckSeries(sample DataSample) []plot.Series
}

// Asc 升序
func Asc[T cmp.Ordered](x, y T) int {
	return cmp.Compare(x, y)
}

// Desc 降序
func Desc[T cmp.Ordered](x, y T) int {
	return Asc(x, y) * -1
}
