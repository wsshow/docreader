package docreader

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// RtfReader 用于读取 .rtf 文件
type RtfReader struct{}

// ReadText 读取 RTF 文件的文本内容（简单提取纯文本）
func (r *RtfReader) ReadText(filePath string) (string, error) {
	// 读取文件内容
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", WrapError("RtfReader.ReadText", filePath, ErrFileRead)
	}

	content := string(data)

	// 简单的 RTF 文本提取
	// 移除 RTF 控制字符
	content = removeRtfControls(content)

	return content, nil
}

// GetMetadata 获取 RTF 文件的元数据
func (r *RtfReader) GetMetadata(filePath string) (map[string]string, error) {
	metadata := make(map[string]string)

	// 获取文件信息
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, WrapError("RtfReader.GetMetadata", filePath, ErrFileNotFound)
	}

	metadata["size"] = fmt.Sprintf("%d", fileInfo.Size())
	metadata["modified"] = fileInfo.ModTime().String()

	return metadata, nil
}

// removeRtfControls 移除 RTF 控制字符，提取纯文本
func removeRtfControls(content string) string {
	// 移除 RTF 头部
	re := regexp.MustCompile(`\\rtf\d+`)
	content = re.ReplaceAllString(content, "")

	// 移除控制字
	re = regexp.MustCompile(`\\[a-z]+\d*\s?`)
	content = re.ReplaceAllString(content, "")

	// 移除花括号
	content = strings.ReplaceAll(content, "{", "")
	content = strings.ReplaceAll(content, "}", "")

	// 移除多余的空白
	re = regexp.MustCompile(`\s+`)
	content = re.ReplaceAllString(content, " ")

	return strings.TrimSpace(content)
}

// ReadWithConfig 根据配置读取 RTF 文件，返回结构化结果
func (r *RtfReader) ReadWithConfig(filePath string, config *ReadConfig) (*DocumentResult, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, WrapError("RtfReader.ReadWithConfig", filePath, ErrFileRead)
	}

	content := string(data)
	content = removeRtfControls(content)
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
