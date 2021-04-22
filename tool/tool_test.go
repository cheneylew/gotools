package tool

import (
	"testing"
	"fmt"
)

type U struct {
	Name string
	Age int
}

func TestArrFile(t *testing.T)  {
	a := []U{U{Name:"1", Age:1},U{Name:"2", Age:2},U{Name:"3", Age:3}}
	fmt.Println(ArrAny(func(a U) bool {return a.Age<2}, a))
}
