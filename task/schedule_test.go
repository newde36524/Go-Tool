package task

import (
	"fmt"
	"testing"
)

func TestSc(t *testing.T) {
	Schedule([]string{"1", "2", "3", "4", "5", "6"}, []func(v string){
		func(v string) {
			fmt.Println("A:", v)
		},
		func(v string) {
			fmt.Println("B:", v)
		},
		func(v string) {
			fmt.Println("C:", v)
		},
	})
}
