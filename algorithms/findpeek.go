package algorithms

import (
	"fmt"
	"strings"
)

// PeeksAndValleys 波峰波谷
type PeeksAndValleys struct {
	Data      []float64 // 原始数据
	Diff      []float64 // 一阶差分
	PosPeak   []int     // 波峰位置存储
	PosValley []int     // 波谷位置存储
	Pcnt      int       // 所识别的波峰计数
	Vcnt      int       // 所识别的波谷计数
}

// InitPV 创建并初始化波峰波谷
func InitPV(data []float64) *PeeksAndValleys {
	n := len(data)
	if n == 0 {
		return nil
	}
	pv := PeeksAndValleys{
		Data:      data,
		Diff:      make([]float64, n),
		PosPeak:   make([]int, n),
		PosValley: make([]int, n),
	}
	for i, v := range data {
		pv.Diff[i] = 0
		pv.PosPeak[i] = -1
		pv.PosValley[i] = -1
		_ = v
	}
	pv.Pcnt = 0
	pv.Vcnt = 0
	return &pv
}

// Find 找波峰波谷
func (this *PeeksAndValleys) Find() {
	n := len(this.Data)
	//step 1 :首先进行前向差分，并归一化
	for i := 0; i < n-1; i++ {
		//int samplei=Sample[i]/1000;
		c := this.Data[i]
		//int samplei1=Sample[i+1]/1000;
		b := this.Data[i+1]
		//printf("%d   %d \n",samplei1,samplei);
		if b-c > 0 {
			this.Diff[i] = 1
		} else if b-c < 0 {
			this.Diff[i] = -1
		} else {
			this.Diff[i] = 0
		}
	}

	//step 2 :对相邻相等的点进行领边坡度处理
	for i := 0; i < n-1; i++ {
		if this.Diff[i] == 0 {
			if i == (n - 2) {
				if this.Diff[i-1] >= 0 {
					this.Diff[i] = 1
				} else {
					this.Diff[i] = -1
				}
			} else {
				if this.Diff[i+1] >= 0 {
					this.Diff[i] = 1
				} else {
					this.Diff[i] = -1
				}
			}
		}
	}
	//step 3 :对相邻相等的点进行领边坡度处理
	for i := 0; i < n-1; i++ {
		//波峰识别
		if this.Diff[i+1]-this.Diff[i] == -2 {
			this.PosPeak[this.Pcnt] = i + 1
			this.Pcnt++
		} else if this.Diff[i+1]-this.Diff[i] == 2 {
			//波谷识别
			this.PosValley[this.Vcnt] = i + 1
			this.Vcnt++
		}
	}
}

func (this *PeeksAndValleys) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("Peak count %d \n", this.Pcnt))
	for _, v := range this.PosPeak {
		sb.WriteString(fmt.Sprintf("-%d", v))
	}
	sb.WriteString(fmt.Sprintf("\nValley count %d \n", this.Vcnt))
	for _, v := range this.PosValley {
		sb.WriteString(fmt.Sprintf("-%d", v))
	}
	return sb.String()
}
