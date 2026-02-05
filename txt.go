package docreader

import (
	"fmt"
	"os"
	"strings"
)

// TxtReader 用于读取 .txt 文件
type TxtReader struct{}

// ReadText 读取 TXT 文件的文本内容
func (r *TxtReader) ReadText(filePath string) (string, error) {
	// 读取文件内容
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", WrapError("TxtReader.ReadText", filePath, ErrFileRead)
	}

	return string(data), nil
}

// GetMetadata 获取 TXT 文件的元数据
func (r *TxtReader) GetMetadata(filePath string) (map[string]string, error) {
	metadata := make(map[string]string)

	// 获取文件信息
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, WrapError("TxtReader.GetMetadata", filePath, ErrFileNotFound)
	}

	metadata["size"] = fmt.Sprintf("%d", fileInfo.Size())
	metadata["modified"] = fileInfo.ModTime().String()

	return metadata, nil
}

// ReadWithConfig 根据配置读取 TXT 文件，返回结构化结果
func (r *TxtReader) ReadWithConfig(filePath string, config *ReadConfig) (*DocumentResult, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, WrapError("TxtReader.ReadWithConfig", filePath, ErrFileRead)
	}

	content := string(data)
	lines := strings.Split(content, "\n")

	result := &DocumentResult{
		FilePath:   filePath,
		TotalPages: 1,
		Pages:      make([]PageContent, 0),
		Metadata:   make(map[string]string),
	}

	// 获取元数据
	metadata, _ := r.GetMetadata(filePath)
	result.Metadata = metadata

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
