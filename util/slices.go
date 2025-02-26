package util

import (
	"fmt"
	"reflect"
	"unsafe"
)

// 对齐到CPU缓存行（通常64字节）
const cacheLineSize = 64

func cacheAlignedSize(elemSize int) int {
	return (elemSize + cacheLineSize - 1) / cacheLineSize * cacheLineSize
}

// 针对AVX-512要求的64字节对齐
const simdAlignment = 64

// TODO 有除零的bug
func simdOptimizedCapacity(length, elemSize int) int {
	elementsPerSIMD := simdAlignment / elemSize
	return (length + elementsPerSIMD - 1) / elementsPerSIMD * elementsPerSIMD
}

const (
	debug = true
)

// ShrinkWithAlignment 带内存对齐的智能缩容函数
func ShrinkWithAlignment[T any](s []T) []T {
	currentLen := len(s)
	currentCap := cap(s)

	// 缩容触发条件：当前容量超过长度的两倍
	if currentCap <= 2*currentLen {
		return s
	}

	// 获取元素类型信息
	elemType := reflect.TypeOf(s).Elem()
	elemSize := int(elemType.Size())
	alignSize := int(elemType.Align())

	// 计算对齐后的新容量
	newCap := calculateAlignedCapacity(currentLen, elemSize, alignSize)

	// 创建对齐的新切片
	newSlice := make([]T, currentLen, newCap)
	copy(newSlice, s)

	// 验证内存对齐
	if debug {
		verifyAlignment(newSlice, elemSize, alignSize)
	}

	return newSlice
}

// 计算对齐容量（核心算法）
func calculateAlignedCapacity(length, elemSize, align int) int {
	minCapacity := length + length/2 // 最小缩容容量
	minCapacity = length
	// 计算满足对齐的最小容量
	alignedCapacity := (minCapacity + align - 1) &^ (align - 1)

	//// 保证至少保留25%的缓冲空间
	//if alignedCapacity < length+length/4 {
	//	alignedCapacity = length + length/4
	//}
	_ = elemSize
	return alignedCapacity
}

// 内存对齐验证函数（调试用）
func verifyAlignment[T any](s []T, elemSize, align int) {
	header := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	addr := header.Data

	if addr%uintptr(align) != 0 {
		panic(fmt.Sprintf("内存未对齐! 地址: %v 对齐要求: %d", addr, align))
	}

	totalSize := cap(s) * elemSize
	if totalSize%align != 0 {
		panic(fmt.Sprintf("容量未对齐! 总字节数: %d 对齐要求: %d",
			totalSize, align))
	}
}
