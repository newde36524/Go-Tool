package geohash

import (
	"fmt"
	"testing"
)

func TestArround(t *testing.T) {
	result := Around(116.3967, 44.9999)
	t.Log(result)
	fmt.Println(result)
}
