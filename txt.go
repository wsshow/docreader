package docreader

import (
	"fmt"
	"os"
)

// TxtReader 用于读取 .txt 文件
type TxtReader struct{}

// ReadText 读取 TXT 文件的文本内容
func (r *TxtReader) ReadText(filePath string) (string, error) {
	// 读取文件内容
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read txt file: %w", err)
	}

	return string(data), nil
}

// GetMetadata 获取 TXT 文件的元数据
func (r *TxtReader) GetMetadata(filePath string) (map[string]string, error) {
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
