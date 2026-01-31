package docreader

import (
	"fmt"

	"github.com/ledongthuc/pdf"
)

// PdfReader 用于读取 .pdf 文件
type PdfReader struct{}

// ReadText 读取 PDF 文件的文本内容
func (r *PdfReader) ReadText(filePath string) (string, error) {
	// 打开 PDF 文件
	f, reader, err := pdf.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open PDF file: %w", err)
	}
	defer f.Close()

	// 获取总页数
	totalPages := reader.NumPage()

	var content string

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

		content += text
		content += "\n\n--- 第 " + fmt.Sprintf("%d", pageNum) + " 页 ---\n\n"
	}

	return content, nil
}

// GetMetadata 获取 PDF 文件的元数据
func (r *PdfReader) GetMetadata(filePath string) (map[string]string, error) {
	f, reader, err := pdf.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open PDF file: %w", err)
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
