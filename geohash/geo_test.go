package geohash

import (
	"fmt"
	"testing"
)

func TestArround(t *testing.T) {
	result := Around2(116.3967, 44.9999, 8)
	t.Error(result)
	fmt.Println(result)
}
