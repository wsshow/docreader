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
		return "", fmt.Errorf("failed to open csv file: %w", err)
	}
	defer file.Close()

	// 创建 CSV 读取器
	reader := csv.NewReader(file)

	// 读取所有记录
	records, err := reader.ReadAll()
	if err != nil {
		return "", fmt.Errorf("failed to read csv content: %w", err)
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
		return nil, fmt.Errorf("failed to open csv file: %w", err)
	}
	defer file.Close()

	// 创建 CSV 读取器
	reader := csv.NewReader(file)

	// 读取所有记录
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read csv content: %w", err)
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
		return nil, fmt.Errorf("failed to open csv file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read csv content: %w", err)
	}

	return records, nil
}
