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
		return "", WrapError("XlsxReader.ReadText", filePath, ErrFileOpen)
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
		return nil, WrapError("XlsxReader.GetMetadata", filePath, ErrFileOpen)
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
		return nil, WrapError("XlsxReader.GetSheetData", filePath, ErrFileOpen)
	}
	defer f.Close()

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, WrapError("XlsxReader.GetSheetData", filePath, ErrSheetNotFound)
	}

	return rows, nil
}

// GetAllSheetsData 获取所有工作表的数据
func (r *XlsxReader) GetAllSheetsData(filePath string) (map[string][][]string, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, WrapError("XlsxReader.GetAllSheetsData", filePath, ErrFileOpen)
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

// ReadWithConfig 根据配置读取 XLSX 文件，返回结构化结果
func (r *XlsxReader) ReadWithConfig(filePath string, config *ReadConfig) (*DocumentResult, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, WrapError("XlsxReader.ReadWithConfig", filePath, ErrFileOpen)
	}
	defer f.Close()

	sheets := f.GetSheetList()
	totalSheets := len(sheets)

	result := &DocumentResult{
		FilePath:   filePath,
		TotalPages: totalSheets,
		Pages:      make([]PageContent, 0),
		Metadata:   make(map[string]string),
	}

	// 获取元数据
	metadata, _ := r.GetMetadata(filePath)
	result.Metadata = metadata

	// 确定要读取的工作表
	var sheetsToRead []int
	sheetNamesSet := make(map[string]bool)

	// 如果指定了工作表名称
	if config != nil && config.SheetNames != nil {
		for _, name := range config.SheetNames {
			sheetNamesSet[name] = true
		}
	}

	// 如果有详细的页面配置
	if config != nil && len(config.PageConfigs) > 0 {
		// 从PageConfigs中提取工作表索引
		for _, pageConfig := range config.PageConfigs {
			if pageConfig.PageIndex >= 0 && pageConfig.PageIndex < totalSheets {
				sheetsToRead = append(sheetsToRead, pageConfig.PageIndex)
			}
		}
	} else if config != nil && (len(config.PageSelector.Indexes) > 0 || len(config.PageSelector.Ranges) > 0) {
		sheetsToRead = determinePagesToRead(config, totalSheets)
	} else if len(sheetNamesSet) > 0 {
		// 根据工作表名称确定索引
		for i, sheetName := range sheets {
			if sheetNamesSet[sheetName] {
				sheetsToRead = append(sheetsToRead, i)
			}
		}
	} else {
		// 读取所有工作表
		sheetsToRead = make([]int, 0, totalSheets)
		for i := 0; i < totalSheets; i++ {
			sheetsToRead = append(sheetsToRead, i)
		}
	}

	// 构建页面行配置映射
	pageLineMap := buildPageLineMap(config, totalSheets)

	var contentBuilder strings.Builder
	totalLines := 0

	for _, sheetIndex := range sheetsToRead {
		if sheetIndex < 0 || sheetIndex >= totalSheets {
			continue
		}

		sheetName := sheets[sheetIndex]
		rows, err := f.GetRows(sheetName)
		if err != nil {
			continue
		}

		// 将每行转换为字符串
		lines := make([]string, 0, len(rows))
		for rowIndex, row := range rows {
			if len(row) == 0 {
				continue
			}

			var lineBuilder strings.Builder
			lineBuilder.WriteString(fmt.Sprintf("Row %d: ", rowIndex))
			for colIndex, cell := range row {
				if colIndex > 0 {
					lineBuilder.WriteString(" | ")
				}
				lineBuilder.WriteString(cell)
			}
			lines = append(lines, lineBuilder.String())
		}

		// 根据配置筛选行
		var filteredLines []string
		if lineConfig, ok := pageLineMap[sheetIndex]; ok {
			filteredLines = filterLinesForPage(lines, lineConfig)
		} else {
			filteredLines = lines
		}

		pageContent := PageContent{
			PageNumber: sheetIndex,
			PageName:   sheetName,
			Lines:      filteredLines,
			TotalLines: len(filteredLines),
		}

		result.Pages = append(result.Pages, pageContent)
		totalLines += len(filteredLines)

		// 构建完整内容
		contentBuilder.WriteString(fmt.Sprintf("\n=== 工作表: %s ===\n\n", sheetName))
		for _, line := range filteredLines {
			contentBuilder.WriteString(line)
			contentBuilder.WriteString("\n")
		}
		contentBuilder.WriteString("\n")
	}

	result.TotalLines = totalLines
	result.Content = contentBuilder.String()

	return result, nil
}
