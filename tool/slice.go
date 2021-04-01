package tool

import (
	"strings"
	"reflect"
	"fmt"
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

func ArrToInterfaces(arr interface{}) ([]interface{},error) {
	v := reflect.ValueOf(arr)
	fmt.Println(v.Kind())
	if v.Kind() != reflect.Slice {
		return nil, ErrorMSG("arr必须为slice类型")
	}
	var res []interface{}
	switch reflect.TypeOf(arr).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(arr)
		for i := 0; i < s.Len(); i++ {
			res = append(res, s.Index(i).Interface())
		}
	}
	return res, nil
}