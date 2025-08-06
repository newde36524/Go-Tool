package slicex

import (
	"context"
	"fmt"
	"strconv"
	"testing"
)

func Test_Slice(t *testing.T) {
	data := map[string]interface{}{"a": 1, "c": 3, "b": 2}
	SortEachMap(data, func(a, b string) bool {
		return a < b
	}, func(k string, v interface{}) {
		fmt.Println(k, v)
	})
}

func Test_Select(t *testing.T) {
	fmt.Println(Select([]int{1, 2, 3}, func(e int) int {
		return e
	}).GenerateWithNum(32))
}

func Test_Generate(t *testing.T) {
	fmt.Println(Generate(3, func(e int) int {
		return e
	}))
}

func Test_Distinct(t *testing.T) {
	type Item struct {
		Id   int64
		Name string
	}
	data := []Item{
		{Id: 1, Name: "1"},
		{Id: 1, Name: "2"},
		{Id: 1, Name: "2"},
		{Id: 1, Name: "3"},
	}
	fmt.Println(data)
	fmt.Println(Distinct(data, func(item Item) int64 {
		return item.Id
	}))
}

func Test_SortEach(t *testing.T) {
	type Item struct {
		Id   int64
		Name string
	}
	data := []Item{
		{Id: 1, Name: "1"},
		{Id: 2, Name: "2"},
		{Id: 3, Name: "2"},
		{Id: 4, Name: "3"},
	}
	fmt.Println(data)
	SortEach(data, func(a, b Item) bool {
		return a.Id > b.Id
	}, func(t Item) bool {
		fmt.Println(t)
		return true
	})
}

func Test_SumFloat(t *testing.T) {
	type Item struct {
		Price float64
	}
	data := []Item{
		{Price: 1.1111},
		{Price: 2.2223},
		{Price: 3.3443},
		{Price: 4.444},
	}
	fmt.Println(data)
	fmt.Println(Select(data, func(in Item) float64 {
		return in.Price
	}).SumFloat())
}

func Test_JoinStr(t *testing.T) {
	type Item struct {
		Data string
	}
	data := []Item{
		{Data: "A"},
		{Data: "B"},
		{Data: "C"},
		{Data: "D"},
	}
	fmt.Println(data)
	fmt.Println(Select(data, func(in Item) string {
		return in.Data
	}).JoinStr(","))
}

func TestSort(t *testing.T) {
	type Data struct {
		Model string
	}
	type Wrapper struct {
		Datas []Data
	}
	wrapper := Wrapper{
		Datas: []Data{
			{
				Model: "1",
			},
			{
				Model: "3",
			},
			{
				Model: "2",
			},
		},
	}
	fmt.Println(1, wrapper.Datas)
	Sort(wrapper.Datas, func(a Data, b Data) bool {
		return a.Model < b.Model
	})
	fmt.Println(2, wrapper.Datas)
}

func Test_MultiProcess(t *testing.T) {
	ch := Select([]int{1, 2, 3}, func(in int) int { return in }).MultiProcess(context.Background(), 1, func(ctx context.Context, in int) (any, bool) {
		return in, true
	})
	for v := range ch {
		fmt.Println(v)
	}
}

func Test_MultiProcess_Wait(t *testing.T) {
	Generate(500, func(index int) int {
		return index
	}).MultiProcess(context.Background(), 10, func(ctx context.Context, in int) (any, bool) {
		fmt.Println(in)
		return in, true
	}).Wait()
}

func Test_MultiProcess_From(t *testing.T) {
	data := []int{1, 2, 3}
	From(data).MultiProcess(context.Background(), 1, func(ctx context.Context, in int) (any, bool) {
		fmt.Println(in)
		return in, true
	}).Wait()
}

func Test_MultiProcess_Fromr(t *testing.T) {
	data := []int{1, 2, 3}
	src := Fromr[int, string](data).MultiProcess(context.Background(), 1, func(ctx context.Context, in int) (string, bool) {
		return strconv.Itoa(in) + "xxx", true
	})
	for v := range src {
		fmt.Println(v)
	}
}

func Test_MergeToMap(t *testing.T) {
	data := []int{1, 2, 3}
	src := MergeToMap(data, func(in int, m map[int]string) {
		m[in] = strconv.Itoa(in) + "xxx"
	})
	for k, v := range src {
		fmt.Println(k, v)
	}
}

func Test_PadLeft(t *testing.T) {
	data := []string{"1", "2", "3"}
	src := PadLeft(data, "0", 12)
	fmt.Println(src.JoinStr(""))
}

func Test_PadRight(t *testing.T) {
	data := []string{"1", "2", "3"}
	src := PadRight(data, "0", 12)
	fmt.Println(src.JoinStr(""))
}
