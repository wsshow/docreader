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
		return "", fmt.Errorf("failed to read rtf file: %w", err)
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
		return nil, fmt.Errorf("failed to get file info: %w", err)
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
