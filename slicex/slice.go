package slicex

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type (
	Slicex[T any]     []T
	Slicexr[T, R any] Slicex[T]
)

// Generate 根据条件生成指定长度的数组
func Generate[T any](num int, callback func(index int) T) Slicex[T] {
	result := make(Slicex[T], num)
	for i := range num {
		result[i] = callback(i)
	}
	return result
}

// GenerateWithNum 生成指定长度的数组并拷贝原数组的数据
func GenerateWithNum[S ~[]T, T any](data S, num int) Slicex[T] {
	source := make(Slicex[T], num)
	for i, v := range data {
		if i < len(source) {
			source[i] = v
		}
	}
	return source
}

// At 获取数组指定下标的值，不存在时返回默认值
func At[S ~[]T, T any](s S, index int) (t T) {
	for i, v := range s {
		if i == index {
			return v
		}
	}
	return
}

// First 获取数组中满足条件的第一个元素，不存在时返回默认值
func First[S ~[]T, T any](s S, callback func(in T) bool) (t T, ok bool) {
	for i := range s {
		if v := s[i]; callback(v) {
			return v, true
		}
	}
	return
}

// Exist 判断数组中是否存在满足条件的元素
func Exist[S ~[]T, T any](s S, callback func(in T) bool) (ok bool) {
	for i := range s {
		if v := s[i]; callback(v) {
			return true
		}
	}
	return
}

// Select 原数组转换成新类型的数组
func Select[S ~[]T, T, R any](data S, callback func(in T) R) Slicex[R] {
	return SelectWithIndex(data, func(_ int, in T) R {
		return callback(in)
	})
}

// Select 原数组转换成新类型的数组
func Contains[S ~[]T, T any](data S, callback func(in T) bool) bool {
	return slices.ContainsFunc(data, callback)
}

// Merge 原数组合并成新类型
func Merge[S ~[]T, T, R any](data S, callback func(in T, r *R)) *R {
	r := new(R)
	for _, v := range data {
		callback(v, r)
	}
	return r
}

// MergeToMap 原数组合并到map
func MergeToMap[S ~[]T, T, R any, K comparable, M map[K]R](data S, callback func(in T, m M)) M {
	m := make(M)
	for _, v := range data {
		callback(v, m)
	}
	return m
}

// Distinct 过滤重复数据返回新数组
func Distinct[S ~[]T, T, R any](data S, condition func(in T) R) Slicex[T] {
	mp := make(map[any]T)
	for _, v := range data {
		key := condition(v)
		if _, ok := mp[key]; !ok {
			mp[key] = v
		}
	}
	var result Slicex[T]
	for _, v := range mp {
		result = append(result, v)
	}
	return result
}

// SelectWithIndex 原数组转换成新类型的数组
func SelectWithIndex[S ~[]T, T, R any](data S, callback func(index int, in T) R) Slicex[R] {
	result := make(Slicex[R], len(data))
	for i := range data {
		result[i] = callback(i, data[i])
	}
	return result
}

// Filter 返回过滤后的数组
func Filter[S ~[]T, T any](data S, callback func(in T) bool) Slicex[T] {
	var result Slicex[T]
	for _, v := range data {
		if callback(v) {
			result = append(result, v)
		}
	}
	return result
}

// FilterR 返回过滤后的数组
func FilterR[S ~[]T, T, R any](data S, callback func(in T) bool) Slicexr[T, R] {
	return Slicexr[T, R](Filter(data, callback))
}

// Repeat 根据指定元素值生成指定长度的数据
func Repeat[T any](v T, count int) Slicex[T] {
	result := make(Slicex[T], count)
	for i := range count {
		result[i] = v
	}
	return result
}

// SortEachMap 顺序遍历map
func SortEachMap[M ~map[K]V, K comparable, V any](data M, cmp func(a, b K) bool, fn func(k K, v V)) {
	var ks Slicex[K]
	for k := range data {
		ks = append(ks, k)
	}
	sort.Slice(ks, func(i, j int) bool { return cmp(ks[i], ks[j]) })
	for _, k := range ks {
		fn(k, data[k])
	}
}

// SortEach 顺序遍历
func SortEach[S ~[]T, T any](data S, cmp func(a, b T) bool, fn func(T) bool) {
	sort.Slice(data, func(i, j int) bool { return cmp(data[i], data[j]) })
	for _, v := range data {
		if !fn(v) {
			break
		}
	}
}

// Sort 数组排序
func Sort[S ~[]T, T comparable](data S, cmp func(a, b T) bool) {
	sort.Slice(data, func(i, j int) bool { return cmp(data[i], data[j]) })
}

// SumFloat 浮点数求和
func SumFloat(data ...float64) (sum float64) {
	var max int
	for _, v := range data {
		str := fmt.Sprintf("%#v", v)
		tmp := len(str) - strings.Index(str, ".") - 1
		if tmp > max {
			max = tmp
		}
	}
	mul, _ := strconv.Atoi("1" + strings.Repeat("0", max))
	num := 0
	for _, v := range data {
		num += int(v * float64(mul))
	}
	sum = float64(num) / float64(mul)
	return
}

// SumInt64 整形求和
func SumInt64(data ...int64) (sum int64) {
	var num int64
	for _, v := range data {
		num += v
	}
	return num
}

// Page 分页获取数组的数据
func Page[S ~[]T, T any](data S, pageNum, pageSize int) (result S) {
	if len(data) <= pageSize {
		return data
	}
	if pageNum < 1 {
		pageNum = 1
	}
	index := (pageNum - 1) * pageSize
	if index > len(data) {
		return
	}
	if index+pageSize > len(data) {
		return data[index:]
	}
	return data[index : index+pageSize]
}

// ToMap 数组转成map
func ToMap[S ~[]T, K comparable, T, V any](data S, getKey func(in T) (k K, v V)) (result map[K]V) {
	result = make(map[K]V, 0)
	for _, v := range data {
		k, v := getKey(v)
		result[k] = v
	}
	return
}

// ToMapSlice 数组转成map,并且数值是slice
func ToMapSlice[S ~[]T, K comparable, T, V any](data S, getKey func(in T) (k K, v V)) (result map[K][]V) {
	result = make(map[K][]V, 0)
	for _, v := range data {
		k, v := getKey(v)
		result[k] = append(result[k], v)
	}
	return
}

// GetKeys 获取map的所有key值
func GetKeys[M ~map[K]V, K comparable, V any](data M) Slicex[K] {
	result := make(Slicex[K], 0)
	for k := range data {
		result = append(result, k)
	}
	return result
}

// HasKey 判断map是否存在key值
func HasKey[M ~map[K]V, K comparable, V any](data M, key K) bool {
	_, ok := data[key]
	return ok
}

// MultiProcessResult 多任务处理返回只读通道
type MultiProcessResult[T any] <-chan T

// Wait 等待多任务处理完成
func (m MultiProcessResult[T]) Wait() {
	for range m { //nolint:revive
	}
}

// Range 处理多任务返回结果
func (m MultiProcessResult[T]) Range(fn func(in T)) {
	for v := range m {
		fn(v)
	}
}

// MultiProcess 多任务并行处理数据且把返回值推送到即将关闭的通道
func MultiProcess[S ~[]T, T, R any](ctx context.Context, workers int, datas S, process func(ctx context.Context, in T) (R, bool)) MultiProcessResult[R] {
	output := make(chan R, len(datas))
	go func(ctx context.Context, output chan R) {
		defer close(output)
		wg, ch := new(sync.WaitGroup), make(chan struct{}, workers)
		for _, v := range datas {
			ch <- struct{}{}
			wg.Add(1)
			go func(ctx context.Context, data T, wg *sync.WaitGroup) {
				if v, ok := process(ctx, data); ok {
					output <- v
				}
				<-ch
				wg.Done()
			}(ctx, v, wg)
		}
		wg.Wait()
	}(ctx, output)
	return output
}

// CopyTo 拷贝原数组的数据到新数组
func CopyTo[S ~[]T, T any](src, dst S) {
	for i, v := range src {
		if i >= len(dst) {
			break
		}
		dst[i] = v
	}
}

// ToJson 数组json序列化
func ToJson[S ~[]T, T any](data S) string {
	bs, _ := json.Marshal(data) //nolint:errchkjson
	return string(bs)
}

// PadLeft 拷贝原数组并往左补充指定元素
func PadLeft[S ~[]T, T any](data S, v T, totalWidth int) Slicex[T] {
	if totalWidth < 0 {
		return Slicex[T]{}
	}
	result, cnt := Repeat(v, totalWidth), len(data)
	if totalWidth-cnt < 0 {
		cnt = totalWidth
	}
	copy(result[totalWidth-cnt:], data[len(data)-cnt:])
	return result
}

// PadRight 拷贝原数组并往右补充指定元素
func PadRight[S ~[]T, T any](data S, v T, totalWidth int) Slicex[T] {
	if totalWidth < 0 {
		return Slicex[T]{}
	}
	result, cnt := Repeat(v, totalWidth), len(data)
	if totalWidth-cnt < 0 {
		cnt = totalWidth
	}
	copy(result[:cnt], data)
	return result
}

// ------------------------ Slicex[T] ------------------------

// From 数据源
func From[S ~[]T, T any](data S) Slicex[T] {
	return Slicex[T](data)
}

// At 获取数组指定下标的值，不存在时返回默认值
func (s Slicex[T]) At(index int) T {
	return At(s, index)
}

// First 获取数组中满足条件的第一个元素，不存在时返回默认值
func (s Slicex[T]) First(callback func(in T) bool) (t T, ok bool) {
	return First(s, callback)
}

// Exist 判断数组中是否存在满足条件的元素
func (s Slicex[T]) Exist(callback func(in T) bool) (ok bool) {
	return Exist(s, callback)
}

// CopyTo 拷贝原数组的数据到新数组
func (s Slicex[T]) CopyTo(dst Slicex[T]) {
	CopyTo(s, dst)
}

// GenerateWithNum 生成指定长度的数组并拷贝原数组的数据
func (s Slicex[T]) GenerateWithNum(num int) Slicex[T] {
	return GenerateWithNum(s, num)
}

// ToJson 数组json序列化
func (s Slicex[T]) ToJson() string {
	return ToJson(s)
}

// MultiProcess 多任务并行处理数据且把返回值推送到即将关闭的通道
func (s Slicex[T]) MultiProcess(ctx context.Context, workers int, process func(ctx context.Context, in T) (any, bool)) MultiProcessResult[any] {
	return MultiProcess(ctx, workers, s, process)
}

// Filter 返回过滤后的数组
func (s Slicex[T]) Filter(callback func(in T) bool) Slicex[T] {
	return Filter(s, callback)
}

// ToTree 数组转树返回根节点和map
func (s Slicex[T]) ToTree(getIdAndParentId func(in T) (selfId, parentId string)) (root []*Node[T], mp Query[T]) {
	return ToTree(s, getIdAndParentId)
}

// SumFloat 浮点数求和
func (s Slicex[T]) SumFloat() float64 {
	var data any = s
	if v, ok := data.(Slicex[float64]); ok {
		return SumFloat(v...)
	}
	return 0
}

// SumInt64 整形求和
func (s Slicex[T]) SumInt64() int64 {
	var data any = s
	if v, ok := data.(Slicex[int64]); ok {
		return SumInt64(v...)
	}
	return 0
}

// JoinStr 合并字符串
func (s Slicex[T]) JoinStr(sep string) string {
	var data any = s
	if v, ok := data.(Slicex[string]); ok {
		return strings.Join(v, sep)
	}
	return ""
}

// PadLeft 生成指定元素的数组，拷贝原数组并往左补充元素
func (s Slicex[T]) PadLeft(t T, length int) Slicex[T] {
	return PadLeft(s, t, length)
}

// PadRight 生成指定元素的数组，拷贝原数组并往右补充元素
func (s Slicex[T]) PadRight(t T, length int) Slicex[T] {
	return PadRight(s, t, length)
}

// PadRight 生成指定元素的数组，拷贝原数组并往右补充元素
func (s Slicex[T]) Raw() []T {
	return s
}

// ------------------------ Slicexr[T, R] ------------------------

// Fromr 数据源
func Fromr[T, R any](data Slicex[T]) Slicexr[T, R] {
	return Slicexr[T, R](data)
}

// Select 原数组转换成新类型的数组
func (s Slicexr[T, R]) Select(callback func(in T) R) Slicex[R] {
	return Select(s, callback)
}

// At 获取数组指定下标的值，不存在时返回默认值
func (s Slicexr[T, R]) At(index int) T {
	return At(s, index)
}

// First 获取数组中满足条件的第一个元素，不存在时返回默认值
func (s Slicexr[T, R]) First(callback func(in T) bool) (t T, ok bool) {
	return First(s, callback)
}

// Exist 判断数组中是否存在满足条件的元素
func (s Slicexr[T, R]) Exist(callback func(in T) bool) (ok bool) {
	return Exist(s, callback)
}

// CopyTo 拷贝原数组的数据到新数组
func (s Slicexr[T, R]) CopyTo(dst Slicexr[T, R]) {
	CopyTo(s, dst)
}

// GenerateWithNum 生成指定长度的数组并拷贝原数组的数据
func (s Slicexr[T, R]) GenerateWithNum(num int) Slicexr[T, R] {
	return Slicexr[T, R](GenerateWithNum(s, num))
}

// ToJson 数组json序列化
func (s Slicexr[T, R]) ToJson() string {
	return ToJson(s)
}

// MultiProcess 多任务并行处理数据且把返回值推送到即将关闭的通道
func (s Slicexr[T, R]) MultiProcess(ctx context.Context, workers int, process func(ctx context.Context, in T) (R, bool)) MultiProcessResult[R] {
	return MultiProcess(ctx, workers, Slicex[T](s), process)
}

// Filter 返回过滤后的数组
func (s Slicexr[T, R]) Filter(callback func(in T) bool) Slicexr[T, R] {
	return Slicexr[T, R](Filter(s, callback))
}

// ToTree 数组转树返回根节点和map
func (s Slicexr[T, R]) ToTree(getIdAndParentId func(in T) (selfId, parentId string)) (root []*Node[T], mp Query[T]) {
	return ToTree(s, getIdAndParentId)
}

// Merge 原数组合并成新类型
func (s Slicexr[T, R]) Merge(callback func(in T, r *R)) *R {
	return Merge(s, callback)
}

// Distinct 过滤重复数据返回新数组
func (s Slicexr[T, R]) Distinct(condition func(in T) R) Slicexr[T, R] {
	return Slicexr[T, R](Distinct(s, condition))
}

// SumFloat 浮点数求和
func (s Slicexr[T, R]) SumFloat() float64 {
	var data any = s
	if v, ok := data.(Slicex[float64]); ok {
		return SumFloat(v...)
	}
	return 0
}

// SumInt64 整形求和
func (s Slicexr[T, R]) SumInt64() int64 {
	var data any = s
	if v, ok := data.(Slicex[int64]); ok {
		return SumInt64(v...)
	}
	return 0
}

// JoinStr 合并字符串
func (s Slicexr[T, R]) JoinStr(sep string) string {
	var data any = s
	if v, ok := data.(Slicex[string]); ok {
		return strings.Join(v, sep)
	}
	return ""
}

// PadLeft 生成指定元素的数组，拷贝原数组并往左补充元素
func (s Slicexr[T, R]) PadLeft(t T, length int) Slicexr[T, R] {
	return Slicexr[T, R](PadLeft(s, t, length))
}

// PadRight 生成指定元素的数组，拷贝原数组并往右补充元素
func (s Slicexr[T, R]) PadRight(t T, length int) Slicexr[T, R] {
	return Slicexr[T, R](PadRight(s, t, length))
}

// SplitN 将  src 分成两个数组，第一个数组长度为 n，第二个数组长度为 src 长度减去 n
func SplitN[S ~[]T, T any](src S, n int) (x, y S) {
	if len(src) <= n {
		return src, y
	}
	x = src[:n]
	y = src[n:]
	return
}
