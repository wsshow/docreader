package docreader

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// 支持的文档格式列表
var supportedFormats = []string{".docx", ".pdf", ".xlsx", ".pptx", ".txt", ".csv", ".md", ".markdown", ".rtf"}

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

// CleanContent 使用默认配置清理文档内容
func (d *Document) CleanContent() {
	d.Content = CleanText(d.Content)
}

// CleanContentWith 使用自定义清理器清理文档内容
func (d *Document) CleanContentWith(cleaner *TextCleaner) {
	d.Content = cleaner.Clean(d.Content)
}

// CleanContentMinimal 使用最小配置清理文档内容
func (d *Document) CleanContentMinimal() {
	d.Content = CleanTextMinimal(d.Content)
}

// CleanContentAggressive 使用激进配置清理文档内容
func (d *Document) CleanContentAggressive() {
	d.Content = CleanTextAggressive(d.Content)
}

// GetSupportedFormats 返回当前支持的文档格式列表
func GetSupportedFormats() []string {
	formats := make([]string, len(supportedFormats))
	copy(formats, supportedFormats)
	return formats
}

// IsFormatSupported 检查指定的文件格式是否被支持
func IsFormatSupported(ext string) bool {
	ext = strings.ToLower(ext)
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	return slices.Contains(supportedFormats, ext)
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

// ReadDocumentWithClean 读取文档并自动应用默认清理
func ReadDocumentWithClean(filePath string) (*Document, error) {
	doc, err := ReadDocument(filePath)
	if err != nil {
		return nil, err
	}
	doc.CleanContent()
	return doc, nil
}

// ReadDocumentWithCleanConfig 读取文档并应用自定义清理配置
func ReadDocumentWithCleanConfig(filePath string, cleaner *TextCleaner) (*Document, error) {
	doc, err := ReadDocument(filePath)
	if err != nil {
		return nil, err
	}
	doc.CleanContentWith(cleaner)
	return doc, nil
}
