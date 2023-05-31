package excel

import (
	"github.com/tealeg/xlsx"
	"os"
)

// 读取excel表格中的数据
// 示例：
//
//	type Person struct {
//		Name  string
//		Phone string
//	}
//
//	excel.ReadData("test.xlsx", func(cells []*xlsx.Cell) {
//			person = append(person, Person{
//				Name:  cells[0].String(),
//				Phone: cells[1].String(),
//			})
//		})
func ReadData(excelPath string, rowFunc func([]*xlsx.Cell)) {
	// 打开Excel文件
	xlFile, err := xlsx.OpenFile(excelPath)
	if err != nil {
		panic(err.Error())
	}

	// 假设文件中只有一个工作表
	sheet := xlFile.Sheets[0]

	// 遍历每一行（忽略表头）
	for _, row := range sheet.Rows[1:] {
		// 读取单元格数据
		rowFunc(row.Cells)
	}

	// 删除Excel文件
	err = os.Remove(excelPath)
	if err != nil {
		panic(err.Error())
	}
}
