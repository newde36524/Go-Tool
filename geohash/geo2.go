package geohash

import "bytes"

const (
	BASE32                = "0123456789bcdefghjkmnpqrstuvwxyz"
	MAX_LATITUDE  float64 = 90
	MIN_LATITUDE  float64 = -90
	MAX_LONGITUDE float64 = 180
	MIN_LONGITUDE float64 = -180
)

var (
	bits   = []int{16, 8, 4, 2, 1}
	base32 = []byte(BASE32)
)

type Box struct {
	MinLat, MaxLat float64 // 纬度
	MinLng, MaxLng float64 // 经度
}

func (this *Box) Width() float64 {
	return this.MaxLng - this.MinLng
}

func (this *Box) Height() float64 {
	return this.MaxLat - this.MinLat
}

// 输入值：纬度，经度，精度(geohash的长度)
// 返回geohash, 以及该点所在的区域
func Encode(latitude, longitude float64, precision int) string {
	var geohash bytes.Buffer
	var minLat, maxLat float64 = MIN_LATITUDE, MAX_LATITUDE
	var minLng, maxLng float64 = MIN_LONGITUDE, MAX_LONGITUDE
	var mid float64 = 0

	bit, ch, length, isEven := 0, 0, 0, true
	for length < precision {
		if isEven {
			if mid = (minLng + maxLng) / 2; mid < longitude {
				ch |= bits[bit]
				minLng = mid
			} else {
				maxLng = mid
			}
		} else {
			if mid = (minLat + maxLat) / 2; mid < latitude {
				ch |= bits[bit]
				minLat = mid
			} else {
				maxLat = mid
			}
		}

		isEven = !isEven
		if bit < 4 {
			bit++
		} else {
			geohash.WriteByte(base32[ch])
			length, bit, ch = length+1, 0, 0
		}
	}

	// b := &Box{
	// 	MinLat: minLat,
	// 	MaxLat: maxLat,
	// 	MinLng: minLng,
	// 	MaxLng: maxLng,
	// }

	return geohash.String()
}

//Around2 计算坐标点geohash
func Around2(lng, lat float64, precision int) []string {
	list := make([]string, 0)
	list = append(list, Encode(lat, lng, precision))
	list = append(list, Encode(lat+latUnit, lng, precision))
	list = append(list, Encode(lat-latUnit, lng, precision))
	list = append(list, Encode(lat, lng+lngUnit, precision))
	list = append(list, Encode(lat, lng-lngUnit, precision))
	list = append(list, Encode(lat+latUnit, lng+lngUnit, precision))
	list = append(list, Encode(lat+latUnit, lng-lngUnit, precision))
	list = append(list, Encode(lat-latUnit, lng+lngUnit, precision))
	list = append(list, Encode(lat-latUnit, lng-lngUnit, precision))
	return list
}
