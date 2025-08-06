package slicex

import (
	"sort"
	"strconv"
)

// 定义一个数字字符串数组的类型
type NumberStrings []string

// Len 实现 sort.Interface 接口的方法
func (ns NumberStrings) Len() int {
	return len(ns)
}

// Less 实现 sort.Interface 接口的方法
// 将字符串转换为整数进行比较
func (ns NumberStrings) Less(i, j int) bool {
	iVal, err := strconv.Atoi(ns[i])
	if err != nil {
		// 如果转换失败，按原字符串比较
		return ns[i] < ns[j]
	}
	jVal, err := strconv.Atoi(ns[j])
	if err != nil {
		return false
	}
	return iVal < jVal
}

// Swap 实现 sort.Interface 接口的方法
func (ns NumberStrings) Swap(i, j int) {
	ns[i], ns[j] = ns[j], ns[i]
}

func SortStrings(data []string) {
	numberStrings := NumberStrings(data)
	sort.Sort(numberStrings)
}
