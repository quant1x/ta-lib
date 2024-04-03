package blas

import "github.com/wcharczuk/go-chart/v2"

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

// WaveTendency 波段趋势
type WaveTendency = int

const (
	TendencyPrice          WaveTendency = 0 // 股价主导
	TendencyLinear         WaveTendency = 1 // 线性趋势主导
	TendencyPriceAndLinear WaveTendency = 2 // 股价和趋势并存
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
	//Match(waves Waves) Pattern
	ExportSeries(sample DataSample) []chart.Series
}
