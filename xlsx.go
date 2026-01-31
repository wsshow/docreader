package docreader

import (
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"
)

// XlsxReader 用于读取 .xlsx 文件
type XlsxReader struct{}

// ReadText 读取 XLSX 文件的文本内容
func (r *XlsxReader) ReadText(filePath string) (string, error) {
	// 打开 Excel 文件
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open xlsx file: %w", err)
	}
	defer f.Close()

	var builder strings.Builder

	// 获取所有工作表
	sheets := f.GetSheetList()

	for _, sheetName := range sheets {
		builder.WriteString(fmt.Sprintf("\n=== 工作表: %s ===\n\n", sheetName))

		// 获取工作表中的所有行
		rows, err := f.GetRows(sheetName)
		if err != nil {
			builder.WriteString(fmt.Sprintf("Failed to read sheet: %v\n", err))
			continue
		}

		// 逐行输出
		for rowIndex, row := range rows {
			// 跳过空行
			if len(row) == 0 {
				continue
			}

			builder.WriteString(fmt.Sprintf("第 %d 行: ", rowIndex+1))

			for colIndex, cell := range row {
				if colIndex > 0 {
					builder.WriteString(" | ")
				}
				builder.WriteString(cell)
			}
			builder.WriteString("\n")
		}
		builder.WriteString("\n")
	}

	return builder.String(), nil
}

// GetMetadata 获取 XLSX 文件的元数据
func (r *XlsxReader) GetMetadata(filePath string) (map[string]string, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open xlsx file: %w", err)
	}
	defer f.Close()

	metadata := make(map[string]string)

	// 获取文档属性
	props, err := f.GetDocProps()
	if err == nil {
		metadata["title"] = props.Title
		metadata["subject"] = props.Subject
		metadata["creator"] = props.Creator
		metadata["description"] = props.Description
		metadata["created"] = props.Created
		metadata["modified"] = props.Modified
		metadata["category"] = props.Category
		metadata["keywords"] = props.Keywords
	}

	// 获取工作表信息
	sheets := f.GetSheetList()
	metadata["sheets"] = strings.Join(sheets, ", ")
	metadata["sheet_count"] = fmt.Sprintf("%d", len(sheets))

	// 获取活动工作表
	activeSheet := f.GetActiveSheetIndex()
	if activeSheet >= 0 && activeSheet < len(sheets) {
		metadata["active_sheet"] = sheets[activeSheet]
	}

	return metadata, nil
}

// GetSheetData 获取指定工作表的结构化数据
func (r *XlsxReader) GetSheetData(filePath, sheetName string) ([][]string, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open xlsx file: %w", err)
	}
	defer f.Close()

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to read sheet: %w", err)
	}

	return rows, nil
}

// GetAllSheetsData 获取所有工作表的数据
func (r *XlsxReader) GetAllSheetsData(filePath string) (map[string][][]string, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open xlsx file: %w", err)
	}
	defer f.Close()

	result := make(map[string][][]string)
	sheets := f.GetSheetList()

	for _, sheetName := range sheets {
		rows, err := f.GetRows(sheetName)
		if err != nil {
			continue
		}
		result[sheetName] = rows
	}

	return result, nil
}
