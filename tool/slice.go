package tool

import (
	"strings"
	"fmt"
	"reflect"
)

//数组字符串包含，即为true
func ArrContain(arr []string, search string) bool {
	for _, item := range arr {
		if strings.Contains(item, search) {
			return true
		}
	}

	return false
}

//数组字符串完全相对，即为true
func ArrIn(arr []string, search string) bool {
	for _, item := range arr {
		if item == search {
			return true
		}
	}

	return false
}

func ArrToInterfaces(slices interface{}) []interface{} {
	v := reflect.ValueOf(slices)
	if v.Kind() != reflect.Ptr {
		Println("123")
	}
	var res []interface{}
	switch slices.(type) {
	case []interface{}:
		fmt.Println("abcd")
	}

	return res
}