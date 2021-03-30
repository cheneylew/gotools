package tool

import (
	"testing"
	"fmt"
)

func TestConv(t *testing.T)  {
	var a []string
	a = append(a, "ok")
	b := ArrToInterfaces(a)
	fmt.Println(b)
	fmt.Println(ErrorMSG("1234"))
}
