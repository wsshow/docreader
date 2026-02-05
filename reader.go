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

// ConfigurableReader 定义了支持配置的文档读取器接口
type ConfigurableReader interface {
	DocumentReader

	// ReadWithConfig 根据配置读取文档，返回结构化结果
	ReadWithConfig(filePath string, config *ReadConfig) (*DocumentResult, error)
}

// Selector 统一的选择器，用于选择页码或行号
type Selector struct {
	// Indexes 离散的索引列表（从0开始）
	// 例如：[]int{0, 2, 5} 表示选择第0、2、5个元素
	Indexes []int

	// Ranges 连续的范围列表，每个范围是 [start, end]（包含起止，从0开始）
	// 例如：[][2]int{{0, 2}, {5, 7}} 表示选择第0-2和第5-7元素
	Ranges [][2]int
}

// PageConfig 单个页面的配置
type PageConfig struct {
	// PageIndex 页码索引（从0开始）
	PageIndex int

	// LineSelector 该页要读取的行选择器
	LineSelector Selector
}

// ReadConfig 读取配置
type ReadConfig struct {
	// PageSelector 页面选择器，指定要读取哪些页
	// 如果为空（Indexes和Ranges都为nil），则读取所有页
	PageSelector Selector

	// LineSelector 全局行选择器，应用到所有选中的页
	// 如果为空，则读取页面的所有行
	LineSelector Selector

	// PageConfigs 页面级配置，为特定页面指定不同的行选择器
	// 如果某页在 PageConfigs 中有配置，则使用该配置，否则使用全局 LineSelector
	PageConfigs []PageConfig

	// SheetNames 对于XLSX文件，指定要读取的工作表名称
	// 如果为nil，则读取所有工作表
	SheetNames []string
}

// PageContent 表示单页/单工作表/单幻灯片的内容
type PageContent struct {
	// PageNumber 页码/工作表索引/幻灯片编号（从0开始）
	PageNumber int

	// PageName 页面名称（对于XLSX是工作表名称，其他格式为空）
	PageName string

	// Lines 该页的所有行内容
	Lines []string

	// TotalLines 该页的总行数
	TotalLines int
}

// DocumentResult 结构化的文档读取结果
type DocumentResult struct {
	// FilePath 文件路径
	FilePath string

	// Pages 所有页面的内容
	Pages []PageContent

	// TotalPages 文档总页数/工作表数/幻灯片数
	TotalPages int

	// TotalLines 所有页面的总行数
	TotalLines int

	// Metadata 文档元数据
	Metadata map[string]string

	// Content 完整的文本内容（所有页面拼接）
	Content string
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

// ReadDocumentWithConfig 根据配置读取文档，返回结构化结果
func ReadDocumentWithConfig(filePath string, config *ReadConfig) (*DocumentResult, error) {
	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, WrapError("ReadDocumentWithConfig", filePath, ErrFileNotFound)
	}

	ext := strings.ToLower(filepath.Ext(filePath))

	var reader ConfigurableReader

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
		return nil, WrapError("ReadDocumentWithConfig", filePath, ErrUnsupportedFormat)
	}

	return reader.ReadWithConfig(filePath, config)
}

// NewReadConfig 创建一个新的读取配置
func NewReadConfig() *ReadConfig {
	return &ReadConfig{}
}

// WithPages 设置要读取的页码（离散索引）
func (c *ReadConfig) WithPages(pages ...int) *ReadConfig {
	c.PageSelector.Indexes = pages
	return c
}

// WithPageRange 设置要读取的页码范围 [start, end]
func (c *ReadConfig) WithPageRange(start, end int) *ReadConfig {
	c.PageSelector.Ranges = append(c.PageSelector.Ranges, [2]int{start, end})
	return c
}

// WithLines 设置要读取的行号（离散索引，应用到所有页）
func (c *ReadConfig) WithLines(lines ...int) *ReadConfig {
	c.LineSelector.Indexes = lines
	return c
}

// WithLineRange 设置要读取的行号范围 [start, end]（应用到所有页）
func (c *ReadConfig) WithLineRange(start, end int) *ReadConfig {
	c.LineSelector.Ranges = append(c.LineSelector.Ranges, [2]int{start, end})
	return c
}

// WithSheetNames 设置要读取的工作表名称（仅用于XLSX）
func (c *ReadConfig) WithSheetNames(names ...string) *ReadConfig {
	c.SheetNames = names
	return c
}

// AddPageConfig 为指定页面添加特定的行选择器
// pageIndex: 页码索引（从0开始）
// lineIndexes: 该页要读取的行号（离散索引）
// lineRanges: 该页要读取的行号范围
func (c *ReadConfig) AddPageConfig(pageIndex int, lineIndexes []int, lineRanges [][2]int) *ReadConfig {
	if c.PageConfigs == nil {
		c.PageConfigs = make([]PageConfig, 0)
	}
	c.PageConfigs = append(c.PageConfigs, PageConfig{
		PageIndex: pageIndex,
		LineSelector: Selector{
			Indexes: lineIndexes,
			Ranges:  lineRanges,
		},
	})
	return c
}

// AddPageLines 为指定页面添加离散行号配置（简化方法）
func (c *ReadConfig) AddPageLines(pageIndex int, lines ...int) *ReadConfig {
	return c.AddPageConfig(pageIndex, lines, nil)
}

// AddPageLineRange 为指定页面添加行范围配置（简化方法）
func (c *ReadConfig) AddPageLineRange(pageIndex int, start, end int) *ReadConfig {
	return c.AddPageConfig(pageIndex, nil, [][2]int{{start, end}})
}
