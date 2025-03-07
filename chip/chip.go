package chip

import (
	"errors"
	"gitee.com/quant1x/engine/datasource/base"
	"gitee.com/quant1x/engine/factors"
	"gitee.com/quant1x/exchange"
	"gitee.com/quant1x/num"
	"math"
	"slices"
	"sort"
	"strings"
	"time"
)

// TechSignal 技术信号位掩码 (支持组合信号)
type TechSignal uint64

const (
	ShortTermRebound  TechSignal = 1 << iota // 短线止跌     (0000 0001)
	ShortTermBreakout                        // 短线突破     (0000 0010)
	VolumeBreakout                           // 放量突破     (0000 0100)
	VolumeBreakdown                          // 放量破位     (0000 1000)
	StrongSupport                            // 强支撑       (0001 0000)
	// 可继续扩展其他信号...
)

// 组合信号示例
var (
	ReboundWithSupport = ShortTermRebound | StrongSupport   // 0001 0001
	BreakoutSignals    = ShortTermBreakout | VolumeBreakout // 0000 0110
)

// Has 判断是否包含某信号
func (ts TechSignal) Has(signal TechSignal) bool {
	return ts&signal != 0
}

// Add 添加信号
func (ts *TechSignal) Add(signal TechSignal) {
	*ts |= signal
}

// Remove 移除信号
func (ts *TechSignal) Remove(signal TechSignal) {
	*ts &^= signal
}

// 转换为可读字符串
func (ts TechSignal) String() string {
	var builder strings.Builder
	if ts.Has(ShortTermRebound) {
		builder.WriteString("短线止跌|")
	}
	if ts.Has(ShortTermBreakout) {
		builder.WriteString("短线突破|")
	}
	if ts.Has(VolumeBreakout) {
		builder.WriteString("放量突破|")
	}
	if ts.Has(VolumeBreakdown) {
		builder.WriteString("放量破位|")
	}
	if ts.Has(StrongSupport) {
		builder.WriteString("强支撑|")
	}
	str := builder.String()
	if len(str) > 0 {
		str = str[:len(str)-1] // 去除末尾的 |
	}
	return str
}

// DailyData 日线数据结构
type DailyData struct {
	base.KLine
	TurnoverRate float64
	Avg          float64
}

// PeakSignal 峰值信号
type PeakSignal struct {
	Closest            float64 // 最近峰值价格
	Extremum           float64 // 极值价格, 筹码最密集的价位
	CurrentToPeakVol   float64 // 当前价格到峰值的累计筹码
	CurrentToPeakRatio float64 // 当前价格到峰值的累计筹码占比 (0~1)
	PeakVolume         float64 // 峰值对应筹码量
	PeakRatio          float64 // 峰值价位本身的筹码占比 (0~1)

	// 以下为扩展分析字段
	LeftVolume    float64 // 峰值左侧筹码总量
	RightVolume   float64 // 峰值右侧筹码总量
	AvgHoldDays   int     // 峰值区平均持仓天数
	Concentration float64 // 筹码集中度 (标准差)

}

// ChipDistribution 筹码分布计算器
type ChipDistribution struct {
	chip         map[float64]float64 // 当前筹码分布
	data         []DailyData         // 日线数据
	config       Config              // 计算配置
	Capital      float64             // 流通股本
	Digits       int                 // 小数点位数
	HoldingPrice float64             // 目标价格
	LastClose    float64             // 最新收盘
	High         float64             // 最高价
	Low          float64             // 最低价
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
		chip:   make(map[float64]float64),
		config: cfg,
		Digits: 2, // 默认2位小数点
	}
}

const (
	yearPeriod = 1
)

// FiveYearsAgoJanFirst 获取五年前的1月1日零点
func FiveYearsAgoJanFirst() time.Time {
	now := time.Now()
	return time.Date(now.Year()-yearPeriod, 1, 1, 0, 0, 0, 0, now.Location())
}

// LoadCSV 加载数据
func (cd *ChipDistribution) LoadCSV(code, date string) error {
	klines := base.CheckoutKLines(code, date)
	f10 := factors.GetL5F10(code)
	if f10 == nil {
		return errors.New("获取F10数据异常")
	}
	cd.Capital = f10.Capital
	cd.Digits = f10.DecimalPoint

	high := math.Inf(-1)
	low := math.Inf(1)
	lastClose := 0.00
	// 计算数据的有效起始日期
	activeDeadline := FiveYearsAgoJanFirst().Format(exchange.TradingDayDateFormat)
	// 预分配切片容量
	cd.data = make([]DailyData, 0, len(klines))
	for _, record := range klines {
		if record.Date < activeDeadline {
			continue
		}
		data := DailyData{
			KLine:        record,
			TurnoverRate: 100 * (record.Volume / cd.Capital),
			Avg:          record.Amount / record.Volume,
		}
		data.Open = num.Decimal(data.Open, cd.Digits)
		data.Close = num.Decimal(data.Close, cd.Digits)
		data.High = num.Decimal(data.High, cd.Digits)
		data.Low = num.Decimal(data.Low, cd.Digits)
		data.Avg = num.Decimal(data.Avg, cd.Digits)
		if data.High > high {
			high = data.High
		}
		if data.Low < low {
			low = data.Low
		}
		cd.data = append(cd.data, data)
		lastClose = data.Close
	}
	cd.data = slices.Clip(cd.data)
	cd.LastClose = lastClose
	if !math.IsInf(high, 0) {
		cd.High = high
	}
	if !math.IsInf(low, 0) {
		cd.Low = low
	}
	return nil
}

func (cd *ChipDistribution) RealVolume(proportion float64) float64 {
	return cd.Capital * proportion
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
	singlePrice := num.Decimal(day.Close, cd.Digits)

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
	priceGrid := generatePriceGrid(day.Low, day.High, cd.config.PriceStep, cd.Digits)
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
	priceGrid := generatePriceGrid(day.Low, day.High, cd.config.PriceStep, cd.Digits)
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
		cd.chip[price] *= 1 - decayRate
	}

	// 合并新筹码
	for price, vol := range newChip {
		cd.chip[price] += vol * decayRate
	}

	// 清理接近零的筹码
	cleanupThreshold := 1e-6 // 根据实际情况调整阈值
	for price := range cd.chip {
		if cd.chip[price] < cleanupThreshold {
			delete(cd.chip, price)
		}
	}

	// 保存当前状态
	//cd.saveChipState(day.Date)
}

// 生成价格区间网格
func generatePriceGrid(low, high, step float64, digits int) []float64 {
	scale := math.Pow10(digits)
	lowInt := int(math.Round(low * scale))
	highInt := int(math.Round(high * scale))
	stepInt := int(math.Round(step * scale))

	if stepInt <= 0 {
		stepInt = 1
	}

	var grid []float64
	for priceInt := lowInt; priceInt <= highInt; priceInt += stepInt {
		price := float64(priceInt) / scale
		price = num.Decimal(price, digits) // 确保四舍五入
		grid = append(grid, price)
	}
	return grid
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

func sortMapKeys(m map[float64]float64) []float64 {
	keys := make([]float64, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Float64s(keys)
	return keys
}

func findMaxPeak(current float64, prices []float64, data map[float64]float64) (price, volume float64) {
	maxVol := 0.0
	var peak float64
	for _, p := range prices {
		if data[p] > maxVol {
			maxVol = data[p]
			peak = p
		}
	}
	high := peak
	low := current
	if high < low {
		high, low = low, high
	}
	var vol float64
	for _, p := range prices {
		if p >= low && p <= high {
			vol += data[p]
		}
	}
	return peak, vol
}

// FindMainPeaks 直接使用当前chip
func (cd *ChipDistribution) FindMainPeaks(targetPrice float64) (upper, lower PeakSignal, err error) {
	if len(cd.chip) == 0 {
		err = errors.New("无可用筹码数据")
		return
	}

	total := cd.calculateTotalVolume(cd.chip)
	if total <= 0 {
		err = errors.New("总筹码量为零")
		return
	}
	cd.HoldingPrice = targetPrice
	sorted := sortMapKeys(cd.chip)
	peaks := cd.findLocalPeaks(sorted, cd.chip)

	// 分离上下峰值
	var lowerPeaks, upperPeaks []float64
	for _, p := range peaks {
		if p < cd.HoldingPrice {
			lowerPeaks = append(lowerPeaks, p)
		} else if p > cd.HoldingPrice {
			upperPeaks = append(upperPeaks, p)
		}
	}

	// 计算特征点
	lower = cd.calculateChipFeature(lowerPeaks, cd.chip, total, false)
	upper = cd.calculateChipFeature(upperPeaks, cd.chip, total, true)
	return
}

// 计算单个特征点信息
func (cd *ChipDistribution) calculateChipFeature(prices []float64, data map[float64]float64, total float64, isUpper bool) PeakSignal {
	var feature PeakSignal

	if len(prices) == 0 {
		return feature
	}

	// 极值计算
	if maxPeakPrice, vol := findMaxPeak(cd.HoldingPrice, prices, data); maxPeakPrice > 0 {
		feature.Extremum = maxPeakPrice
		feature.PeakVolume = data[maxPeakPrice]
		feature.PeakRatio = feature.PeakVolume / total
		feature.CurrentToPeakVol = vol
		feature.CurrentToPeakRatio = feature.CurrentToPeakVol / total
	}

	// 最近峰值计算
	var closestPrice float64
	if isUpper {
		//closestPrice = findMinPrice(prices)
		closestPrice = num.Min2(prices)
	} else {
		//closestPrice = findMaxPrice(prices)
		closestPrice = num.Max2(prices)
	}

	if closestPrice > 0 {
		feature.Closest = closestPrice
		// 如果极值未找到，使用最近峰值的量能
		if feature.PeakVolume == 0 {
			feature.PeakVolume = data[closestPrice]
			feature.PeakRatio = feature.PeakVolume / total
		}
	}

	return feature
}

// 计算总筹码量
func (cd *ChipDistribution) calculateTotalVolume(data map[float64]float64) float64 {
	total := 0.0
	for _, vol := range data {
		total += vol
	}
	return total
}
