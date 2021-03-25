package excel

import (
	"strings"
	"fmt"
	"path"
	"os"
	"math/rand"
	"time"
	"github.com/cheneylew/gotools/tool"
	"errors"
	"github.com/tealeg/xlsx/v2"
)

type RowItem map[string]interface{}

type SheetData struct {
	SheetName string `json:"sheetName"`
	Data [][]string `json:"data"`
}

type ExcelData struct {
	Sheets []SheetData `json:"sheets"`
}

type DataVerify struct {
	ColIndex int
	Name string
	Items []string
}

func GetSheetItems(sh *xlsx.Sheet) [][]string {
	maxCol := sh.MaxCol
	maxRow := sh.MaxRow
	var data [][]string
	for i:=0; i<maxRow; i++ {
		var row []string
		for j:=0; j<maxCol; j++ {
			cell := sh.Cell(i, j)
			val := cell.String()
			row = append(row, val)
		}
		data = append(data, row)
	}
	return  data
}

func GetSheetHeader(sh *xlsx.Sheet) []string {
	maxCol := sh.MaxCol
	var data []string
	if sh.MaxRow == 0 {
		return data
	}
	for j:=0; j<maxCol; j++ {
		cell := sh.Cell(0, j)
		val := cell.String()
		data = append(data, val)
	}
	return  data
}

func findIndex(arr []string,search string) int {
	for key, value := range arr {
		if value == search {
			return key
		}
	}
	return -1
}

//https://github.com/tealeg/xlsx/blob/master/tutorial/tutorial.adoc
func ReadExcel(file string) (*ExcelData, error) {
	isHttpUrl := strings.HasPrefix(file, "http")
	var wb *xlsx.File
	var err error
	if isHttpUrl {
		bytes, e := tool.DownloadFileToBytes(file)
		if e != nil {
			return nil, errors.New("下载Excel不存在或无法下载："+file+fmt.Sprintf("%v", e))
		}

		wb, err = xlsx.OpenBinary(bytes)
	} else {
		wb, err = xlsx.OpenFile(file)
	}

	if err != nil {
		return nil, tool.ErrorLineWithMSG("Excel已损坏，请重新上传！错误为：%v", err)
	}

	var sheets []SheetData
	for _, sh := range wb.Sheets {
		sh, ok := wb.Sheet[sh.Name]
		if !ok {
			fmt.Println("Sheet does not exist")
			break
		}
		data := GetSheetItems(sh)
		sheets = append(sheets, SheetData{
			Data:data,
			SheetName:sh.Name,
		})
	}
	data := &ExcelData {
		Sheets:sheets,
	}

	return data, nil
}

func ReadTPLExcelHeader(tplFile string) ([]string,map[string]int, error) {
	tplExcel, _ := ReadExcel(tplFile)
	if len(tplExcel.Sheets) == 0 || len(tplExcel.Sheets[0].Data) < 2 {
		return nil, nil, errors.New("Excel不是一个正常的模板!不能写入")
	}
	tplRows := tplExcel.Sheets[0].Data
	tplFields := tplRows[1:2][0]
	tplFieldsMap := make(map[string]int)
	for key, value := range tplFields {
		tplFieldsMap[value] = key
	}

	return  tplFields,tplFieldsMap, nil
}

//https://github.com/tealeg/xlsx/blob/master/tutorial/tutorial.adoc
func WriteExcel(tplFile string, data [][]string) (string, error) {
	fileName := path.Base(tplFile)
	fileExt := path.Ext(tplFile)
	filenameOnly := strings.TrimSuffix(fileName, fileExt)

	isHttpUrl := strings.HasPrefix(tplFile, "http")
	var wb *xlsx.File
	var err error
	if isHttpUrl {
		bytes, e := tool.DownloadFileToBytes(tplFile)
		if e != nil {
			return "", tool.ErrorLineWithMSG("下载Excel TPL模板失败:%v", err)
		}

		wb, err = xlsx.OpenBinary(bytes)
	} else {
		wb, err = xlsx.OpenFile(tplFile)
	}

	if err != nil {
		panic(err)
	}
	sheet := wb.Sheets[0]
	fmt.Println(sheet.MaxRow)
	if len(wb.Sheets) > 1 {
		dataVerifySheet := wb.Sheets[len(wb.Sheets)-1]
		if dataVerifySheet.Name == "数据验证" {
			header := GetSheetHeader(sheet)
			dvData := GetSheetItems(dataVerifySheet)
			if len(dvData) > 0 {
				var dataVerifys []*DataVerify
				for rowIndex, Item := range dvData {
					for colIndex, value := range Item {
						if rowIndex == 0 {
							dataVerify := new(DataVerify)
							dataVerifys = append(dataVerifys, dataVerify)
							dataVerify.Name = value
							dataVerify.ColIndex = findIndex(header, value)
						} else {
							if value != "" {
								dataVerifys[colIndex].Items = append(dataVerifys[colIndex].Items, value)
							}
						}
					}
				}
				if dataVerifys != nil {
					for _, dv := range dataVerifys {
						dd := xlsx.NewDataValidation(sheet.MaxRow, dv.ColIndex, sheet.MaxRow+len(data)+100*1000,  dv.ColIndex, true)
						err = dd.SetDropList(dv.Items)
						if err == nil {
							sheet.AddDataValidation(dd)
						}
					}
				}
			}
		}
	}


	var fields []string
	fieldsMap := make(map[string]int, 0)
	for i:=0; i< sheet.MaxCol ; i++ {
		cell := sheet.Cell(1, i)
		val := cell.String()
		fields = append(fields, val)
		fieldsMap[val] = i
	}

	//config
	fontColor := "00000000"
	fontSize := 10
	borderColor := "00000000"
	borderType := "thin"

	startRowIndex := sheet.MaxRow
	for index := startRowIndex; index < startRowIndex+len(data); index++ {
		for colIndex:=0; colIndex<sheet.MaxCol; colIndex++ {
			cell := sheet.Cell(index, colIndex)
			//https://github.com/tealeg/xlsx/blob/master/tutorial/tutorial.adoc#assigning-a-style
			//fmt.Println(*cell.GetStyle())

			cell.GetStyle().Fill.FgColor = fontColor
			cell.GetStyle().Font.Size = fontSize
			cell.GetStyle().Border.Bottom = borderType
			cell.GetStyle().Border.Right = borderType
			cell.GetStyle().Border.Left = borderType
			cell.GetStyle().Border.BottomColor = borderColor
			cell.GetStyle().Border.RightColor = borderColor
			cell.GetStyle().Border.LeftColor = borderColor

			if colIndex < len(data[index-startRowIndex]) {
				val := data[index-startRowIndex][colIndex]
				cell.SetString(val)
			}
		}
	}

	tmpDirName := "tmp"
	if _, err := os.Stat(tmpDirName); os.IsNotExist(err) {
		err = os.Mkdir(tmpDirName, os.ModePerm)
		return "", tool.ErrorLineWithMSG("写Excel，建临时目录失败!", err)
	}

	rand.Seed(time.Now().UnixNano())
	newFileName := fmt.Sprintf("%s/%s%s", tmpDirName, filenameOnly+fmt.Sprintf("_%v", rand.Intn(1000)), fileExt)

	wb.Save(newFileName)

	return newFileName, nil
}

func WriteExcelWithRowItems(tplFile string, rowItems []*RowItem) (filePath string, err error) {
	tplExcel, err  := ReadExcel(tplFile)
	if err != nil {
		return "", err
	}
	if len(tplExcel.Sheets) == 0 || len(tplExcel.Sheets[0].Data) < 2 {
		return "", errors.New("Excel不是一个正常的模板!不能写入")
	}
	tplRows := tplExcel.Sheets[0].Data
	tplFields := tplRows[1:2][0]
	tplFieldsMap := make(map[string]int)
	for key, value := range tplFields {
		tplFieldsMap[value] = key
	}
	tplRows = tplRows[2:]
	for _, rowItem := range rowItems {
		var tplRow []string
		for i:=0; i< len(tplFields) ; i++ {
			tplRow = append(tplRow, "")
		}
		for key, value := range *rowItem {
			index, ok := tplFieldsMap[key]
			if ok {
				tplRow[index] = fmt.Sprintf("%v", value)
			}
		}
		tplRows = append(tplRows, tplRow)
	}

	p, err := WriteExcel(tplFile, tplRows)
	return p, err
}
