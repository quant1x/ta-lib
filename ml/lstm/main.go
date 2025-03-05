package main

import (
	"fmt"
	"math"
	"math/rand"
)

// 初始化随机种子
func init() {
	rand.Seed(0)
}

// 激活函数
func sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}

func sigmoidDerivative(values []float64) []float64 {
	deriv := make([]float64, len(values))
	for i, v := range values {
		deriv[i] = v * (1 - v)
	}
	return deriv
}

func tanhDerivative(values []float64) []float64 {
	deriv := make([]float64, len(values))
	for i, v := range values {
		deriv[i] = 1 - v*v
	}
	return deriv
}

// 矩阵/向量工具函数
func randMatrix(a, b float64, rows, cols int) [][]float64 {
	matrix := make([][]float64, rows)
	for i := range matrix {
		matrix[i] = make([]float64, cols)
		for j := range matrix[i] {
			matrix[i][j] = rand.Float64()*(b-a) + a
		}
	}
	return matrix
}

func zerosMatrix(rows, cols int) [][]float64 {
	matrix := make([][]float64, rows)
	for i := range matrix {
		matrix[i] = make([]float64, cols)
	}
	return matrix
}

func zerosVector(size int) []float64 {
	return make([]float64, size)
}

func vectorAdd(a, b []float64) []float64 {
	if len(a) != len(b) {
		panic("向量长度不一致")
	}
	result := make([]float64, len(a))
	for i := range result {
		result[i] = a[i] + b[i]
	}
	return result
}

func elementMultiply(a, b []float64) []float64 {
	if len(a) != len(b) {
		panic("向量长度不一致")
	}
	result := make([]float64, len(a))
	for i := range result {
		result[i] = a[i] * b[i]
	}
	return result
}

// LstmParam LSTM参数结构体
type LstmParam struct {
	memCellCt                      int
	xDim                           int
	wg, wi, wf, wo                 [][]float64
	bg, bi, bf, bo                 []float64
	wgDiff, wiDiff, wfDiff, woDiff [][]float64
	bgDiff, biDiff, bfDiff, boDiff []float64
}

func NewLstmParam(memCellCt, xDim int) *LstmParam {
	concatLen := xDim + memCellCt
	return &LstmParam{
		memCellCt: memCellCt,
		xDim:      xDim,
		wg:        randMatrix(-0.1, 0.1, memCellCt, concatLen),
		wi:        randMatrix(-0.1, 0.1, memCellCt, concatLen),
		wf:        randMatrix(-0.1, 0.1, memCellCt, concatLen),
		wo:        randMatrix(-0.1, 0.1, memCellCt, concatLen),
		bg:        randVector(-0.1, 0.1, memCellCt),
		bi:        randVector(-0.1, 0.1, memCellCt),
		bf:        randVector(-0.1, 0.1, memCellCt),
		bo:        randVector(-0.1, 0.1, memCellCt),
		wgDiff:    zerosMatrix(memCellCt, concatLen),
		wiDiff:    zerosMatrix(memCellCt, concatLen),
		wfDiff:    zerosMatrix(memCellCt, concatLen),
		woDiff:    zerosMatrix(memCellCt, concatLen),
		bgDiff:    zerosVector(memCellCt),
		biDiff:    zerosVector(memCellCt),
		bfDiff:    zerosVector(memCellCt),
		boDiff:    zerosVector(memCellCt),
	}
}

func randVector(a, b float64, size int) []float64 {
	vec := make([]float64, size)
	for i := range vec {
		vec[i] = rand.Float64()*(b-a) + a
	}
	return vec
}

// LstmState LSTM状态结构体
type LstmState struct {
	g, i, f, o               []float64
	s, h                     []float64
	bottomDiffH, bottomDiffS []float64
}

func NewLstmState(memCellCt int) *LstmState {
	return &LstmState{
		g:           zerosVector(memCellCt),
		i:           zerosVector(memCellCt),
		f:           zerosVector(memCellCt),
		o:           zerosVector(memCellCt),
		s:           zerosVector(memCellCt),
		h:           zerosVector(memCellCt),
		bottomDiffH: zerosVector(memCellCt),
		bottomDiffS: zerosVector(memCellCt),
	}
}

// LstmNode LSTM节点结构体
type LstmNode struct {
	state *LstmState
	param *LstmParam
	xc    []float64
	sPrev []float64
	hPrev []float64
}

func NewLstmNode(param *LstmParam, state *LstmState) *LstmNode {
	return &LstmNode{
		state: state,
		param: param,
		sPrev: zerosVector(param.memCellCt),
		hPrev: zerosVector(param.memCellCt),
	}
}

// 前向传播时应显式检查输入维度
func (node *LstmNode) bottomDataIs(x []float64, sPrev, hPrev []float64) {
	if len(x) != node.param.xDim {
		panic(fmt.Sprintf("输入维度错误: 预期=%d, 实际=%d",
			node.param.xDim, len(x)))
	}
	if sPrev != nil {
		node.sPrev = sPrev
	}
	if hPrev != nil {
		node.hPrev = hPrev
	}

	// 拼接输入
	xc := make([]float64, len(x)+len(hPrev))
	copy(xc, x)
	copy(xc[len(x):], hPrev)
	node.xc = xc

	// 计算门状态
	node.state.g = node.activate(node.param.wg, node.param.bg, math.Tanh)
	node.state.i = node.activate(node.param.wi, node.param.bi, sigmoid)
	node.state.f = node.activate(node.param.wf, node.param.bf, sigmoid)
	node.state.o = node.activate(node.param.wo, node.param.bo, sigmoid)

	// 更新状态
	node.state.s = vectorAdd(
		elementMultiply(node.state.g, node.state.i),
		elementMultiply(node.sPrev, node.state.f),
	)
	node.state.h = elementMultiply(node.state.s, node.state.o)
}

func (node *LstmNode) activate(weights [][]float64, bias []float64, actFunc func(float64) float64) []float64 {
	result := make([]float64, len(weights))
	for i := range weights {
		sum := bias[i]
		for j := range node.xc {
			sum += weights[i][j] * node.xc[j]
		}
		result[i] = actFunc(sum)
	}
	return result
}

// LstmNetwork LSTM网络结构体
type LstmNetwork struct {
	param      *LstmParam
	nodes      []*LstmNode
	xList      [][]float64
	denseLayer *DenseLayer
}

func NewLstmNetwork(param *LstmParam) *LstmNetwork {
	return &LstmNetwork{
		param:      param,
		denseLayer: NewDenseLayer(param.memCellCt, 1),
	}
}

// DenseLayer 全连接层
type DenseLayer struct {
	Weights [][]float64
	Bias    []float64
}

func NewDenseLayer(inputSize, outputSize int) *DenseLayer {
	return &DenseLayer{
		Weights: randMatrix(-0.1, 0.1, outputSize, inputSize),
		Bias:    randVector(-0.1, 0.1, outputSize),
	}
}

func (d *DenseLayer) Forward(input []float64) []float64 {
	output := make([]float64, len(d.Weights))
	for i := range d.Weights {
		sum := d.Bias[i]
		for j := range input {
			sum += d.Weights[i][j] * input[j]
		}
		output[i] = sum
	}
	return output
}

// MSELoss 损失函数
type MSELoss struct{}

func (m MSELoss) Loss(pred, target []float64) float64 {
	if len(pred) != len(target) {
		panic(fmt.Sprintf("维度不匹配: pred=%d, target=%d", len(pred), len(target)))
	}
	sum := 0.0
	for i := range pred {
		diff := pred[i] - target[i]
		sum += diff * diff
	}
	return sum / 2
}

func (m MSELoss) BottomDiff(pred, target []float64) []float64 {
	diff := make([]float64, len(pred))
	for i := range diff {
		diff[i] = pred[i] - target[i]
	}
	return diff
}

func (param *LstmParam) applyDiff(lr float64) {
	update := func(weights, diffs [][]float64) {
		for i := range weights {
			for j := range weights[i] {
				weights[i][j] -= lr * diffs[i][j]
			}
		}
	}

	update(param.wg, param.wgDiff)
	update(param.wi, param.wiDiff)
	update(param.wf, param.wfDiff)
	update(param.wo, param.woDiff)

	for i := range param.bg {
		param.bg[i] -= lr * param.bgDiff[i]
		param.bi[i] -= lr * param.biDiff[i]
		param.bf[i] -= lr * param.bfDiff[i]
		param.bo[i] -= lr * param.boDiff[i]
	}

	// 重置梯度
	param.resetDiffs()
}

func (param *LstmParam) resetDiffs() {
	param.wgDiff = zerosMatrix(param.memCellCt, param.xDim+param.memCellCt)
	param.wiDiff = zerosMatrix(param.memCellCt, param.xDim+param.memCellCt)
	param.wfDiff = zerosMatrix(param.memCellCt, param.xDim+param.memCellCt)
	param.woDiff = zerosMatrix(param.memCellCt, param.xDim+param.memCellCt)
	param.bgDiff = zerosVector(param.memCellCt)
	param.biDiff = zerosVector(param.memCellCt)
	param.bfDiff = zerosVector(param.memCellCt)
	param.boDiff = zerosVector(param.memCellCt)
}

func (network *LstmNetwork) calculateLoss(targets [][]float64, lossLayer MSELoss, lr float64) float64 {
	totalLoss := 0.0
	seqLen := len(network.nodes)
	if seqLen == 0 {
		return 0.0
	}

	// 处理最后一个节点
	lastNode := network.nodes[seqLen-1]
	h := lastNode.state.h
	pred := network.denseLayer.Forward(h)
	target := targets[seqLen-1]

	// 计算损失和梯度
	loss := lossLayer.Loss(pred, target)
	outputGrad := lossLayer.BottomDiff(pred, target)

	// 全连接层反向传播
	gradInput := make([]float64, network.param.memCellCt)
	for i := range network.denseLayer.Weights[0] {
		gradInput[i] = network.denseLayer.Weights[0][i] * outputGrad[0]
	}

	// 更新全连接层参数
	for i := range network.denseLayer.Weights[0] {
		network.denseLayer.Weights[0][i] -= lr * outputGrad[0] * h[i]
	}
	network.denseLayer.Bias[0] -= lr * outputGrad[0]

	// LSTM反向传播
	diffS := zerosVector(network.param.memCellCt)
	lastNode.topDiffIs(gradInput, diffS)
	totalLoss += loss

	// 处理中间节点
	for idx := seqLen - 2; idx >= 0; idx-- {
		node := network.nodes[idx]
		nextNode := network.nodes[idx+1]

		// 累加梯度
		diffH := vectorAdd(node.state.bottomDiffH, nextNode.state.bottomDiffH)
		diffS := nextNode.state.bottomDiffS

		node.topDiffIs(diffH, diffS)
	}

	return totalLoss
}

func (node *LstmNode) topDiffIs(topDiffH, topDiffS []float64) {
	ds := vectorAdd(elementMultiply(node.state.o, topDiffH), topDiffS)
	do := elementMultiply(node.state.s, topDiffH)
	di := elementMultiply(node.state.g, ds)
	dg := elementMultiply(node.state.i, ds)
	df := elementMultiply(node.sPrev, ds)

	diInput := elementMultiply(sigmoidDerivative(node.state.i), di)
	dfInput := elementMultiply(sigmoidDerivative(node.state.f), df)
	doInput := elementMultiply(sigmoidDerivative(node.state.o), do)
	dgInput := elementMultiply(tanhDerivative(node.state.g), dg)

	// 更新参数梯度
	node.updateParamGradients(diInput, dfInput, doInput, dgInput)

	// 计算bottom diff
	dxc := zerosVector(len(node.xc))
	for i := range node.param.wg {
		for j := range node.xc {
			dxc[j] += node.param.wg[i][j]*dgInput[i] +
				node.param.wi[i][j]*diInput[i] +
				node.param.wf[i][j]*dfInput[i] +
				node.param.wo[i][j]*doInput[i]
		}
	}

	// 保存梯度
	copy(node.state.bottomDiffH, dxc[node.param.xDim:]) // 正确使用xDim作为索引
	copy(node.state.bottomDiffS, elementMultiply(ds, node.state.f))
}

func (node *LstmNode) updateParamGradients(di, df, do, dg []float64) {
	for i := range node.param.wg {
		for j := range node.xc {
			node.param.wgDiff[i][j] += dg[i] * node.xc[j]
			node.param.wiDiff[i][j] += di[i] * node.xc[j]
			node.param.wfDiff[i][j] += df[i] * node.xc[j]
			node.param.woDiff[i][j] += do[i] * node.xc[j]
		}
		node.param.bgDiff[i] += dg[i]
		node.param.biDiff[i] += di[i]
		node.param.bfDiff[i] += df[i]
		node.param.boDiff[i] += do[i]
	}
}

// 主函数和训练循环
func main() {
	memCellCt := 4
	xDim := 1
	seqLength := 20
	epochs := 10
	learningRate := 0.01

	// 生成训练数据
	sinWave := make([][]float64, seqLength)
	targets := make([][]float64, seqLength)
	for i := 0; i < seqLength; i++ {
		sinWave[i] = []float64{math.Sin(float64(i) * 0.1)}
		targets[i] = []float64{math.Sin(float64(i+1) * 0.1)}
	}

	param := NewLstmParam(memCellCt, xDim)
	network := NewLstmNetwork(param)
	lossLayer := MSELoss{}
	fmt.Printf("parmas-init = %+v\n", param)
	for epoch := 0; epoch < epochs; epoch++ {
		totalLoss := 0.0

		// 前向传播
		for t := 0; t < seqLength; t++ {
			network.xList = append(network.xList, sinWave[t])
			if len(network.xList) > len(network.nodes) {
				state := NewLstmState(param.memCellCt)
				node := NewLstmNode(param, state)
				network.nodes = append(network.nodes, node)
			}

			idx := len(network.xList) - 1
			if idx == 0 {
				network.nodes[idx].bottomDataIs(sinWave[idx], nil, nil)
			} else {
				prevNode := network.nodes[idx-1]
				network.nodes[idx].bottomDataIs(sinWave[idx], prevNode.state.s, prevNode.state.h)
			}
		}

		// 反向传播和参数更新
		totalLoss = network.calculateLoss(targets, lossLayer, learningRate)
		param.applyDiff(learningRate)

		// 重置状态
		network.xList = nil
		network.nodes = nil

		if epoch%10 == 0 {
			fmt.Printf("Epoch %d, Loss: %.4f\n", epoch, totalLoss)
		}

	}
	fmt.Printf("parmas finished= %+v\n", param)
}
