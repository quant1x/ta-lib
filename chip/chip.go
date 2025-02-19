package chip

import (
	"errors"
	"fmt"
	"gitee.com/quant1x/engine/datasource/base"
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/num"
	"math"
	"sort"
)

// DailyData 日线数据结构
type DailyData struct {
	base.KLine
	TurnoverRate float64
	Avg          float64
}

type ChipSignle struct {
	Extremum float64
	Closest  float64
}

// ChipDistribution 筹码分布计算器
type ChipDistribution struct {
	chip        map[float64]float64            // 当前筹码分布
	chipHistory map[string]map[float64]float64 // 历史筹码分布
	data        []DailyData                    // 日线数据
	config      Config                         // 计算配置
}

// Config 计算配置参数
type Config struct {
	PriceStep    float64 // 价格最小变动单位
	DecayFactor  float64 // 衰减系数
	ModelType    int     // 计算模型 (1:三角形 2:均匀)
	SearchWindow int     // 峰值搜索窗口大小
}

var (
	defaultConfig = Config{
		PriceStep:    0.01,
		DecayFactor:  0.95,
		ModelType:    1,
		SearchWindow: 5,
	}
)

// NewChipDistribution 创建新的筹码分布计算实例
func NewChipDistribution(cfg Config) *ChipDistribution {
	if cfg.PriceStep <= 0 {
		cfg.PriceStep = 0.01
	}
	if cfg.SearchWindow <= 0 {
		cfg.SearchWindow = 3
	}

	return &ChipDistribution{
		chip:        make(map[float64]float64),
		chipHistory: make(map[string]map[float64]float64),
		config:      cfg,
	}
}

// 加载数据
func (cd *ChipDistribution) LoadCSV(code, date string) error {
	klines := base.CheckoutKLines(code, date)
	f10 := factors.GetL5F10(code)
	capital := f10.Capital
	for _, record := range klines {

		data := DailyData{
			KLine:        record,
			TurnoverRate: record.Volume / capital,
			Avg:          record.Amount / record.Volume,
		}
		cd.data = append(cd.data, data)
	}
	return nil
}

// Calculate 执行筹码分布计算
func (cd *ChipDistribution) Calculate() error {
	if len(cd.data) == 0 {
		return errors.New("无有效数据")
	}

	for _, day := range cd.data {
		switch cd.config.ModelType {
		case 1:
			if err := cd.calculateTriangular(day); err != nil {
				return err
			}
		case 2:
			if err := cd.calculateUniform(day); err != nil {
				return err
			}
		default:
			return errors.New("无效的计算模型类型")
		}
	}
	return nil
}

// 三角形分布计算
func (cd *ChipDistribution) calculateTriangular(day DailyData) error {
	if day.High < day.Low {
		return fmt.Errorf("无效价格区间: 日期 %s", day.Date)
	}

	priceGrid := generatePriceGrid(day.Low, day.High, cd.config.PriceStep)
	tmpChip := make(map[float64]float64)
	h := 2.0 / (day.High - day.Low)

	for _, price := range priceGrid {
		x1 := price
		x2 := price + cd.config.PriceStep
		var area float64

		if price < day.Avg {
			y1 := h / (day.Avg - day.Low) * (x1 - day.Low)
			y2 := h / (day.Avg - day.Low) * (x2 - day.Low)
			area = cd.config.PriceStep * (y1 + y2) / 2
		} else {
			y1 := h / (day.High - day.Avg) * (day.High - x1)
			y2 := h / (day.High - day.Avg) * (day.High - x2)
			area = cd.config.PriceStep * (y1 + y2) / 2
		}

		tmpChip[price] = area * day.Volume
	}

	// 应用衰减和合并
	cd.applyDecayAndMerge(day, tmpChip)
	return nil
}

// 均匀分布计算
func (cd *ChipDistribution) calculateUniform(day DailyData) error {
	priceGrid := generatePriceGrid(day.Low, day.High, cd.config.PriceStep)
	eachVol := day.Volume / float64(len(priceGrid))
	tmpChip := make(map[float64]float64)

	for _, price := range priceGrid {
		tmpChip[price] = eachVol
	}

	cd.applyDecayAndMerge(day, tmpChip)
	return nil
}

// 应用衰减并合并筹码
func (cd *ChipDistribution) applyDecayAndMerge(day DailyData, newChip map[float64]float64) {
	decayRate := day.TurnoverRate / 100 * cd.config.DecayFactor

	// 衰减现有筹码
	for price := range cd.chip {
		cd.chip[price] *= (1 - decayRate)
	}

	// 合并新筹码
	for price, vol := range newChip {
		cd.chip[price] += vol * decayRate
	}

	// 保存当前状态
	cd.saveChipState(day.Date)
}

// 生成价格区间网格
func generatePriceGrid(low, high, step float64) []float64 {
	var grid []float64
	for price := low; price <= high; price += step {
		grid = append(grid, round(price, 2))
	}
	return grid
}

// 保存筹码状态
func (cd *ChipDistribution) saveChipState(date string) {
	currentState := make(map[float64]float64)
	for k, v := range cd.chip {
		currentState[k] = v
	}
	cd.chipHistory[date] = currentState
}

// FindMainPeaks 查找主要筹码峰
func (cd *ChipDistribution) FindMainPeaks(targetPrice float64) (upper, lower ChipSignle, err error) {
	latest, err := cd.getLatestDistribution()
	if err != nil {
		return
	}

	sorted := sortMapKeys(latest)
	peaks := cd.findLocalPeaks(sorted, latest)

	// 分离前后峰值
	var lowerPeaks, upperPeaks []float64
	for _, p := range peaks {
		if p < targetPrice {
			lowerPeaks = append(lowerPeaks, p)
		} else if p > targetPrice {
			upperPeaks = append(upperPeaks, p)
		}
	}

	// 获取最大峰值
	if maxLower := findNextPeak(lowerPeaks, latest); maxLower > 0 {
		lower.Closest = maxLower
	}
	if maxLower := v1findMaxPeak(lowerPeaks, latest); maxLower > 0 {
		lower.Extremum = maxLower
	}
	if maxUpper := findPrevPeak(upperPeaks, latest); maxUpper > 0 {
		upper.Closest = maxUpper
	}
	if maxUpper := v1findMaxPeak(upperPeaks, latest); maxUpper > 0 {
		upper.Extremum = maxUpper
	}

	return upper, lower, nil
}

// 辅助函数：查找局部峰值
func (cd *ChipDistribution) findLocalPeaks(prices []float64, data map[float64]float64) []float64 {
	var peaks []float64
	n := len(prices)
	if n < 3 {
		return nil
	}

	windowSize := cd.config.SearchWindow
	for i := range prices {
		left := max(0, i-windowSize/2)
		right := min(n-1, i+windowSize/2)

		isPeak := true
		currentVol := data[prices[i]]
		for j := left; j <= right; j++ {
			if j == i {
				continue
			}
			if currentVol <= data[prices[j]] {
				isPeak = false
				break
			}
		}

		if isPeak {
			peaks = append(peaks, prices[i])
		}
	}
	return peaks
}

// 辅助函数：获取最新筹码分布
func (cd *ChipDistribution) getLatestDistribution() (map[float64]float64, error) {
	if len(cd.chipHistory) == 0 {
		return nil, errors.New("无可用筹码数据")
	}

	var latestDate string
	for date := range cd.chipHistory {
		if date > latestDate {
			latestDate = date
		}
	}
	return cd.chipHistory[latestDate], nil
}

// 工具函数
func round(val float64, precision int) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func sortMapKeys(m map[float64]float64) []float64 {
	keys := make([]float64, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Float64s(keys)
	return keys
}

func findPrevPeak(prices []float64, data map[float64]float64) float64 {
	n := len(prices)
	if n < 1 {
		num.Float64NaN()
	}
	return prices[0]
}

func findNextPeak(prices []float64, data map[float64]float64) float64 {
	return v3findMaxPeak(prices, data)
}

func v1findMaxPeak(prices []float64, data map[float64]float64) float64 {
	maxVol := 0.0
	var peak float64
	for _, p := range prices {
		if data[p] > maxVol {
			maxVol = data[p]
			peak = p
		}
	}
	return peak
}

func v2findMaxPeak(prices []float64, data map[float64]float64) float64 {
	var peak float64
	n := len(prices)
	for i := 1; i < n-1; i++ {
		prev := prices[i-1]
		curr := prices[i]
		next := prices[i+1]
		if data[curr] > data[prev] && data[curr] > data[next] {
			peak = curr
			fmt.Printf("\t%.2f\n", peak)
			//break
		}
	}
	last := prices[n-1]
	if data[last] > data[peak] {
		peak = data[last]
	}
	return peak
}

func v3findMaxPeak(prices []float64, data map[float64]float64) float64 {
	n := len(prices)
	if n < 1 {
		return num.Float64NaN()
	}
	return prices[n-1]
}
