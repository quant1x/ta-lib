package waves

import (
	"gitee.com/quant1x/num"
)

type Data struct {
	Index   int
	Current float64
	High    float64
	Low     float64
}

// Kind 浪型
type Kind = int8

const (
	Unknown Kind = 0  // 位置状态
	Drive   Kind = 1  // 驱动浪
	Adjust  Kind = -1 // 调整浪
)

// Trend 趋势
type Trend struct {
	Peak    num.DataPoint // 波峰
	Valley  num.DataPoint // 波谷
	Periods int           // 周期数
	State   Kind          // 状态
	Inner   []Trend       // 嵌套浪
}
