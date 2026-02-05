package docreader

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

// CsvReader 用于读取 .csv 文件
type CsvReader struct{}

// ReadText 读取 CSV 文件的文本内容
func (r *CsvReader) ReadText(filePath string) (string, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", WrapError("CsvReader.ReadText", filePath, ErrFileOpen)
	}
	defer file.Close()

	// 创建 CSV 读取器
	reader := csv.NewReader(file)

	// 读取所有记录
	records, err := reader.ReadAll()
	if err != nil {
		return "", WrapError("CsvReader.ReadText", filePath, ErrFileRead)
	}

	var builder strings.Builder

	// 格式化输出
	for rowIndex, record := range records {
		builder.WriteString(fmt.Sprintf("Row %d: ", rowIndex+1))
		builder.WriteString(strings.Join(record, " | "))
		builder.WriteString("\n")
	}

	return builder.String(), nil
}

// GetMetadata 获取 CSV 文件的元数据
func (r *CsvReader) GetMetadata(filePath string) (map[string]string, error) {
	metadata := make(map[string]string)

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, WrapError("CsvReader.GetMetadata", filePath, ErrFileOpen)
	}
	defer file.Close()

	// 创建 CSV 读取器
	reader := csv.NewReader(file)

	// 读取所有记录
	records, err := reader.ReadAll()
	if err != nil {
		return nil, WrapError("CsvReader.GetMetadata", filePath, ErrFileRead)
	}

	metadata["rows"] = fmt.Sprintf("%d", len(records))
	if len(records) > 0 {
		metadata["columns"] = fmt.Sprintf("%d", len(records[0]))
	}

	// 获取文件信息
	fileInfo, err := os.Stat(filePath)
	if err == nil {
		metadata["size"] = fmt.Sprintf("%d", fileInfo.Size())
		metadata["modified"] = fileInfo.ModTime().String()
	}

	return metadata, nil
}

// GetRecords 获取 CSV 文件的结构化数据
func (r *CsvReader) GetRecords(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, WrapError("CsvReader.GetRecords", filePath, ErrFileOpen)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, WrapError("CsvReader.GetRecords", filePath, ErrFileRead)
	}

	return records, nil
}

// ReadWithConfig 根据配置读取 CSV 文件，返回结构化结果
func (r *CsvReader) ReadWithConfig(filePath string, config *ReadConfig) (*DocumentResult, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, WrapError("CsvReader.ReadWithConfig", filePath, ErrFileOpen)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, WrapError("CsvReader.ReadWithConfig", filePath, ErrFileRead)
	}

	result := &DocumentResult{
		FilePath:   filePath,
		TotalPages: 1,
		Pages:      make([]PageContent, 0),
		Metadata:   make(map[string]string),
	}

	// 获取元数据
	metadata, _ := r.GetMetadata(filePath)
	result.Metadata = metadata

	// 将每行记录转换为字符串
	lines := make([]string, 0, len(records))
	for rowIndex, record := range records {
		line := fmt.Sprintf("Row %d: %s", rowIndex+1, strings.Join(record, " | "))
		lines = append(lines, line)
	}

	// 根据配置筛选行
	filteredLines := filterLinesForSinglePage(lines, config)

	pageContent := PageContent{
		PageNumber: 0,
		Lines:      filteredLines,
		TotalLines: len(filteredLines),
	}

	result.Pages = append(result.Pages, pageContent)
	result.TotalLines = len(filteredLines)
	result.Content = strings.Join(filteredLines, "\n")

	return result, nil
}
