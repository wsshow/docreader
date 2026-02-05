package docreader

import (
	"fmt"
	"strings"

	"github.com/ledongthuc/pdf"
)

// PdfReader 用于读取 .pdf 文件
type PdfReader struct{}

// ReadText 读取 PDF 文件的文本内容
func (r *PdfReader) ReadText(filePath string) (string, error) {
	// 打开 PDF 文件
	f, reader, err := pdf.Open(filePath)
	if err != nil {
		return "", WrapError("PdfReader.ReadText", filePath, ErrFileOpen)
	}
	defer f.Close()

	// 获取总页数
	totalPages := reader.NumPage()

	var content strings.Builder

	// 逐页读取文本
	for pageNum := 1; pageNum <= totalPages; pageNum++ {
		page := reader.Page(pageNum)
		if page.V.IsNull() {
			continue
		}

		text, err := page.GetPlainText(nil)
		if err != nil {
			// 如果某页读取失败，继续读取下一页
			continue
		}

		content.WriteString(text)
		content.WriteString("\n\n--- 第 " + fmt.Sprintf("%d", pageNum) + " 页 ---\n\n")
	}

	return content.String(), nil
}

// GetMetadata 获取 PDF 文件的元数据
func (r *PdfReader) GetMetadata(filePath string) (map[string]string, error) {
	f, reader, err := pdf.Open(filePath)
	if err != nil {
		return nil, WrapError("PdfReader.GetMetadata", filePath, ErrFileOpen)
	}
	defer f.Close()

	metadata := make(map[string]string)

	// 获取基本信息
	if !reader.Trailer().IsNull() && !reader.Trailer().Key("Info").IsNull() {
		info := reader.Trailer().Key("Info")

		if title := info.Key("Title"); !title.IsNull() {
			metadata["title"] = title.String()
		}
		if author := info.Key("Author"); !author.IsNull() {
			metadata["author"] = author.String()
		}
		if subject := info.Key("Subject"); !subject.IsNull() {
			metadata["subject"] = subject.String()
		}
		if creator := info.Key("Creator"); !creator.IsNull() {
			metadata["creator"] = creator.String()
		}
		if producer := info.Key("Producer"); !producer.IsNull() {
			metadata["producer"] = producer.String()
		}
		if creationDate := info.Key("CreationDate"); !creationDate.IsNull() {
			metadata["creation_date"] = creationDate.String()
		}
		if modDate := info.Key("ModDate"); !modDate.IsNull() {
			metadata["modification_date"] = modDate.String()
		}
	}

	metadata["pages"] = fmt.Sprintf("%d", reader.NumPage())

	return metadata, nil
}

// ReadWithConfig 根据配置读取 PDF 文件，返回结构化结果
func (r *PdfReader) ReadWithConfig(filePath string, config *ReadConfig) (*DocumentResult, error) {
	f, reader, err := pdf.Open(filePath)
	if err != nil {
		return nil, WrapError("PdfReader.ReadWithConfig", filePath, ErrFileOpen)
	}
	defer f.Close()

	totalPages := reader.NumPage()
	result := &DocumentResult{
		FilePath:   filePath,
		TotalPages: totalPages,
		Pages:      make([]PageContent, 0),
		Metadata:   make(map[string]string),
	}

	// 获取元数据
	metadata, _ := r.GetMetadata(filePath)
	result.Metadata = metadata

	// 确定要读取的页码和每页的行配置
	pageLineMap := buildPageLineMap(config, totalPages)

	var contentBuilder strings.Builder
	totalLines := 0

	// 按页码顺序处理
	for pageIndex := 0; pageIndex < totalPages; pageIndex++ {
		lineConfig, shouldRead := pageLineMap[pageIndex]
		if !shouldRead {
			continue
		}

		// PDF库的页码从1开始，所以需要+1
		page := reader.Page(pageIndex + 1)
		if page.V.IsNull() {
			continue
		}

		text, err := page.GetPlainText(nil)
		if err != nil {
			continue
		}

		// 按行分割
		lines := strings.Split(text, "\n")

		// 根据该页的配置筛选行
		filteredLines := filterLinesForPage(lines, lineConfig)

		pageContent := PageContent{
			PageNumber: pageIndex,
			Lines:      filteredLines,
			TotalLines: len(filteredLines),
		}

		result.Pages = append(result.Pages, pageContent)
		totalLines += len(filteredLines)

		// 构建完整内容
		for _, line := range filteredLines {
			contentBuilder.WriteString(line)
			contentBuilder.WriteString("\n")
		}
		contentBuilder.WriteString(fmt.Sprintf("\n--- 第 %d 页 ---\n\n", pageIndex))
	}

	result.TotalLines = totalLines
	result.Content = contentBuilder.String()

	return result, nil
}
