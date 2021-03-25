package excel

import (
	"testing"
)

func TestHello(t *testing.T)  {
	var rowItems []*RowItem
	rowItems = append(rowItems, &RowItem{
		"sku_code":"A0001",
	})
	WriteExcelWithRowItems("/Users/dejunliu/Desktop/tpl1.xlsx", rowItems)
}