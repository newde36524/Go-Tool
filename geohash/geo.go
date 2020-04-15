package geohash

import (
	"fmt"
	"strings"
)

var (
	max_Lat      float64 = 90
	min_Lat      float64 = -90
	max_Lng      float64 = 180
	min_Lng      float64 = -180
	length               = 20
	latUnit              = (max_Lat - min_Lat) / (1 << 20)
	lngUnit              = (max_Lng - min_Lng) / (1 << 20)
	base32Lookup         = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "b",
		"c", "d", "e", "f", "g", "h", "j", "k", "m", "n", "p",
		"q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
)

//convert .
func convert(min, max, value float64) (list []string) {
	var cvt func(min, max, value float64)
	cvt = func(min, max, value float64) {
		if len(list) > (length - 1) {
			return
		}
		mid := (max + min) / 2
		if value < mid {
			list = append(list, "0")
			cvt(min, mid, value)
		} else {
			list = append(list, "1")
			cvt(mid, max, value)
		}
	}
	cvt(min, max, value)
	fmt.Println(list)
	return
}

//base32Encode .
func base32Encode(str string) string {
	result := make([]string, 0, length)
	for i := 0; i < len(str) && len(str) >= i+5; i += 5 {
		result = append(result, base32Lookup[convertToIndex(str[i:i+5])])
	}
	return strings.Join(result, "")
}

//convertToIndex .
func convertToIndex(str string) int {
	result := 0
	for i := 0; i < len(str); i++ {
		if str[i] == '0' {
			result += 0
		} else {
			result += 1 << (len(str) - 1 - i)
		}
	}
	return result
}

//encode .
func encode(lng, lat float64) string {
	latList := convert(min_Lat, max_Lat, lat)
	lngList := convert(min_Lng, max_Lng, lng)
	sb := make([]string, 0, length)
	for index := 0; index < len(latList); index++ {
		sb = append(sb, lngList[index])
		sb = append(sb, latList[index])
	}
	return base32Encode(strings.Join(sb, ""))
}

//Around 计算坐标点geohash
func Around(lng, lat float64) []string {
	list := make([]string, 0)
	list = append(list, encode(lat, lng))
	list = append(list, encode(lat+latUnit, lng))
	list = append(list, encode(lat-latUnit, lng))
	list = append(list, encode(lat, lng+lngUnit))
	list = append(list, encode(lat, lng-lngUnit))
	list = append(list, encode(lat+latUnit, lng+lngUnit))
	list = append(list, encode(lat+latUnit, lng-lngUnit))
	list = append(list, encode(lat-latUnit, lng+lngUnit))
	list = append(list, encode(lat-latUnit, lng-lngUnit))
	return list
}
