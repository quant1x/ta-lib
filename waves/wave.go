package waves

import "gitee.com/quant1x/num"

type Waves struct {
	Trends     []Trend
	DataLength int
	State      Kind
	LastHigh   num.DataPoint
	LastLow    num.DataPoint
}

func (w *Waves) Len() int {
	return len(w.Trends)
}

func (w *Waves) Push(data Data) {
	high := num.DataPoint{X: data.Index, Y: data.High}
	low := num.DataPoint{X: data.Index, Y: data.Low}
	if w.DataLength == 0 {
		t := Trend{
			State:   Unknown,
			Peak:    high,
			Valley:  low,
			Periods: 1,
		}
		w.Trends = append(w.Trends, t)
	} else {
		n := w.Len()
		lastTrend := &w.Trends[n-1]
		diffHigh := high.Y - w.LastHigh.Y
		diffLow := low.Y - w.LastLow.Y
		switch w.State {
		case Unknown: // 初始状态
			if diffHigh > 0 {
				w.State = Drive
			} else if diffLow < 0 {
				w.State = Adjust
			}
		case Drive: // 升高
			if diffHigh > 0 {
				lastTrend.Peak = high
			} else {
				// 趋势改变
				w.State = Adjust
				lastTrend.Valley = low
				t := Trend{
					State:   Adjust,
					Peak:    high,
					Valley:  low,
					Periods: 1,
				}
				w.Trends = append(w.Trends, t)
			}
		case Adjust: // 降低
			if diffLow < 0 {
				lastTrend.Valley = low
			} else {
				// 趋势改变
				w.State = Drive
				lastTrend.Peak = high
				t := Trend{
					State:   Drive,
					Peak:    high,
					Valley:  low,
					Periods: 1,
				}
				w.Trends = append(w.Trends, t)
			}
		}

	}
	w.LastHigh = high
	w.LastLow = low
	w.DataLength++

}
