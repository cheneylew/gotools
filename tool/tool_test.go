package tool

import (
	"testing"
	"fmt"
)

func TestArrFile(t *testing.T)  {
	a := []int{2,3,4,5,67}
	out := ArrFilter(a, func(item interface{}) bool {
		return item.(int)>10
	})
	fmt.Println(out)
}
