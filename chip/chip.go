package chip

import (
	"errors"
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

// ChipSignle 增强版筹码特征点
type ChipSignle struct {
	Extremum   float64 // 极值价格
	Closest    float64 // 最近峰值价格
	Volume     float64 // 对应筹码量（股）
	Proportion float64 // 占总筹码比例（0-1）
}

// ChipDistribution 筹码分布计算器
type ChipDistribution struct {
	chip        map[float64]float64            // 当前筹码分布
	chipHistory map[string]map[float64]float64 // 历史筹码分布
	data        []DailyData                    // 日线数据
	config      Config                         // 计算配置
	capital     float64                        // 流通股本
	digits      int                            // 小数点位数
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
		DecayFactor:  0.9995,
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
		digits:      2, // 默认2位
	}
}

// LoadCSV 加载数据
func (cd *ChipDistribution) LoadCSV(code, date string) error {
	klines := base.CheckoutKLines(code, date)
	f10 := factors.GetL5F10(code)
	if f10 == nil {
		return errors.New("获取F10数据异常")
	}
	cd.capital = f10.Capital
	cd.digits = f10.DecimalPoint
	for _, record := range klines {
		data := DailyData{
			KLine:        record,
			TurnoverRate: 100 * (record.Volume / cd.capital),
			Avg:          record.Amount / record.Volume,
		}
		data.Open = num.Decimal(data.Open, cd.digits)
		data.Close = num.Decimal(data.Close, cd.digits)
		data.High = num.Decimal(data.High, cd.digits)
		data.Low = num.Decimal(data.Low, cd.digits)
		data.Avg = num.Decimal(data.Avg, cd.digits)
		cd.data = append(cd.data, data)
	}
	return nil
}

func (cd *ChipDistribution) RealVolume(proportion float64) float64 {
	return cd.capital * proportion
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

// 处理单一价格日的特殊逻辑
func (cd *ChipDistribution) handleSinglePriceDay(day DailyData) {
	// 生成唯一价格点（容差处理避免浮点误差）
	//singlePrice  := round(day.Close, 2)
	singlePrice := num.Decimal(day.Close, cd.digits)

	// 构造全量筹码分布
	tmpChip := map[float64]float64{
		singlePrice: day.Volume,
	}

	// 应用衰减合并（需特殊处理衰减率）
	adjustedDay := day
	//adjustedDay.TurnoverRate = 100 // 强制全量换手
	cd.applyDecayAndMerge(adjustedDay, tmpChip)
}

// 三角形分布计算
func (cd *ChipDistribution) calculateTriangular(day DailyData) error {
	// 情况1：处理最高价等于最低价的情况（一字涨停/跌停）
	if day.High == day.Low {
		// 直接全部分配到唯一价格点
		cd.handleSinglePriceDay(day)
		return nil
	}

	// 情况2：常规价格区间校验
	if day.High < day.Low {
		//return fmt.Errorf("无效价格区间: 日期 %s (High=%.2f < Low=%.2f)", day.Date, day.High, day.Low)
		return nil
	}

	// 生成价格网格（包含容差处理）
	priceGrid := generatePriceGrid(day.Low, day.High, cd.config.PriceStep, cd.digits)
	tmpChip := make(map[float64]float64, len(priceGrid)) // 预分配内存优化

	// 计算归一化系数（处理可能的零除问题）
	priceRange := day.High - day.Low
	h := 2.0 / priceRange // 保证概率密度积分为1

	for _, price := range priceGrid {
		x1 := price
		x2 := price + cd.config.PriceStep
		var area float64

		// 分情况处理三角形分布
		if price < day.Avg {
			// 左三角形处理（包含Avg=Low的边界情况）
			// 形态特征: 筹码密度从最低价向均价递增（类似直角三角形）
			// 市场意义: 当日低位买盘活跃，形成下方支撑带
			denominator := day.Avg - day.Low
			if denominator <= 1e-8 { // 处理浮点精度误差
				// 当Avg=Low时退化为矩形分布
				area = cd.config.PriceStep * h
			} else {
				y1 := h / denominator * (x1 - day.Low)
				y2 := h / denominator * (x2 - day.Low)
				area = cd.config.PriceStep * (y1 + y2) / 2
			}
		} else {
			// 右三角形处理（包含Avg=High的边界情况）
			// 形态特征: 筹码密度从均价向最高价递减
			// 市场意义: 当日高位抛压显现，形成上方阻力带
			denominator := day.High - day.Avg
			if denominator <= 1e-8 {
				// 当Avg=High时退化为矩形分布
				area = cd.config.PriceStep * h
			} else {
				y1 := h / denominator * (day.High - x1)
				y2 := h / denominator * (day.High - x2)
				area = cd.config.PriceStep * (y1 + y2) / 2
			}
		}

		tmpChip[price] = area * day.Volume // 面积映射到实际成交量
	}

	// 应用衰减和合并
	cd.applyDecayAndMerge(day, tmpChip)
	return nil
}

// 均匀分布计算
func (cd *ChipDistribution) calculateUniform(day DailyData) error {
	priceGrid := generatePriceGrid(day.Low, day.High, cd.config.PriceStep, cd.digits)
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
	// 在applyDecayAndMerge中自动处理：
	decayRate := day.TurnoverRate / 100 * cd.config.DecayFactor
	decayRate = math.Min(decayRate, 1.0)

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
//
//	TODO: 输出的网格会不会有重复的价格
func generatePriceGrid(low, high, step float64, digits int) []float64 {
	var grid []float64
	for price := low; price <= high; price += step {
		price = num.Decimal(price, digits)
		grid = append(grid, price)
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

func sortMapKeys(m map[float64]float64) []float64 {
	keys := make([]float64, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Float64s(keys)
	return keys
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

// FindMainPeaks 增强版查找主峰（包含筹码量计算）
func (cd *ChipDistribution) FindMainPeaks(targetPrice float64) (upper, lower ChipSignle, err error) {
	latest, err := cd.getLatestDistribution()
	if err != nil {
		return
	}

	// 计算总筹码量
	total := cd.calculateTotalVolume(latest)
	if total <= 0 {
		err = errors.New("总筹码量为零")
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

	// 获取特征点并计算量能
	lower = cd.calculateChipFeature(lowerPeaks, latest, total, false)
	upper = cd.calculateChipFeature(upperPeaks, latest, total, true)

	return upper, lower, nil
}

// 计算单个特征点信息
func (cd *ChipDistribution) calculateChipFeature(prices []float64, data map[float64]float64, total float64, isUpper bool) ChipSignle {
	var feature ChipSignle

	if len(prices) == 0 {
		return feature
	}

	// 极值计算
	if maxPeak := v1findMaxPeak(prices, data); maxPeak > 0 {
		feature.Extremum = maxPeak
		feature.Volume = data[maxPeak]
		feature.Proportion = data[maxPeak] / total
	}

	// 最近峰值计算
	var closestPrice float64
	if isUpper {
		closestPrice = findMinPrice(prices)
	} else {
		closestPrice = findMaxPrice(prices)
	}

	if closestPrice > 0 {
		feature.Closest = closestPrice
		// 如果极值未找到，使用最近峰值的量能
		if feature.Volume == 0 {
			feature.Volume = data[closestPrice]
			feature.Proportion = data[closestPrice] / total
		}
	}

	return feature
}

// 工具函数：查找最小价格（用于上方最近峰）
func findMinPrice(prices []float64) float64 {
	if len(prices) == 0 {
		return 0
	}
	min := prices[0]
	for _, p := range prices {
		if p < min {
			min = p
		}
	}
	return min
}

// 工具函数：查找最大价格（用于下方最近峰）
func findMaxPrice(prices []float64) float64 {
	if len(prices) == 0 {
		return 0
	}
	max := prices[0]
	for _, p := range prices {
		if p > max {
			max = p
		}
	}
	return max
}

// 计算总筹码量
func (cd *ChipDistribution) calculateTotalVolume(data map[float64]float64) float64 {
	total := 0.0
	for _, vol := range data {
		total += vol
	}
	return total
}
