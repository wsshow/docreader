package docreader

import (
	"os"
	"path/filepath"
	"strings"
)

// DocumentReader 定义了文档读取器的通用接口
type DocumentReader interface {
	// ReadText 读取文档的文本内容
	ReadText(filePath string) (string, error)

	// GetMetadata 获取文档元数据
	GetMetadata(filePath string) (map[string]string, error)
}

// Document 表示一个文档及其内容
type Document struct {
	FilePath string
	Content  string
	Metadata map[string]string
}

// ReadDocument 根据文件扩展名自动选择合适的读取器
func ReadDocument(filePath string) (*Document, error) {
	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, WrapError("ReadDocument", filePath, ErrFileNotFound)
	}

	ext := strings.ToLower(filepath.Ext(filePath))

	var reader DocumentReader

	switch ext {
	case ".docx":
		reader = &DocxReader{}
	case ".pdf":
		reader = &PdfReader{}
	case ".xlsx":
		reader = &XlsxReader{}
	case ".pptx":
		reader = &PptxReader{}
	case ".txt":
		reader = &TxtReader{}
	case ".csv":
		reader = &CsvReader{}
	case ".md", ".markdown":
		reader = &MdReader{}
	case ".rtf":
		reader = &RtfReader{}
	default:
		return nil, WrapError("ReadDocument", filePath, ErrUnsupportedFormat)
	}

	content, err := reader.ReadText(filePath)
	if err != nil {
		return nil, err
	}

	metadata, err := reader.GetMetadata(filePath)
	if err != nil {
		metadata = make(map[string]string)
	}

	return &Document{
		FilePath: filePath,
		Content:  content,
		Metadata: metadata,
	}, nil
}
