package docreader

import (
	"archive/zip"
	"encoding/xml"
	"io"
	"strings"
)

// DocxReader 用于读取 .docx 文件
type DocxReader struct{}

// WordDocument 表示 Word 文档的 XML 结构
type WordDocument struct {
	XMLName xml.Name `xml:"document"`
	Body    struct {
		Paragraphs []struct {
			Runs []struct {
				Text string `xml:"t"`
			} `xml:"r"`
		} `xml:"p"`
		Tables []struct {
			Rows []struct {
				Cells []struct {
					Paragraphs []struct {
						Runs []struct {
							Text string `xml:"t"`
						} `xml:"r"`
					} `xml:"p"`
				} `xml:"tc"`
			} `xml:"tr"`
		} `xml:"tbl"`
	} `xml:"body"`
}

// CoreProperties 表示文档核心属性
type CoreProperties struct {
	XMLName     xml.Name `xml:"coreProperties"`
	Title       string   `xml:"title"`
	Subject     string   `xml:"subject"`
	Creator     string   `xml:"creator"`
	Description string   `xml:"description"`
	Created     string   `xml:"created"`
	Modified    string   `xml:"modified"`
}

// ReadText 读取 DOCX 文件的文本内容
func (r *DocxReader) ReadText(filePath string) (string, error) {
	// 打开 zip 文件
	zipReader, err := zip.OpenReader(filePath)
	if err != nil {
		return "", WrapError("DocxReader.ReadText", filePath, ErrFileOpen)
	}
	defer zipReader.Close()

	// 查找并读取 document.xml
	var documentXML []byte
	for _, file := range zipReader.File {
		if file.Name == "word/document.xml" {
			rc, err := file.Open()
			if err != nil {
				return "", WrapError("DocxReader.ReadText", filePath, ErrFileRead)
			}
			documentXML, err = io.ReadAll(rc)
			rc.Close()
			if err != nil {
				return "", WrapError("DocxReader.ReadText", filePath, ErrFileRead)
			}
			break
		}
	}

	if documentXML == nil {
		return "", WrapError("DocxReader.ReadText", filePath, ErrInvalidFormat)
	}

	// 解析 XML
	var doc WordDocument
	if err := xml.Unmarshal(documentXML, &doc); err != nil {
		return "", WrapError("DocxReader.ReadText", filePath, ErrFileParse)
	}

	// 提取文本
	var builder strings.Builder

	// 提取段落文本
	for _, para := range doc.Body.Paragraphs {
		for _, run := range para.Runs {
			builder.WriteString(run.Text)
		}
		builder.WriteString("\n")
	}

	// 提取表格文本
	for _, table := range doc.Body.Tables {
		for _, row := range table.Rows {
			for _, cell := range row.Cells {
				for _, para := range cell.Paragraphs {
					for _, run := range para.Runs {
						builder.WriteString(run.Text)
						builder.WriteString(" ")
					}
				}
				builder.WriteString("\t")
			}
			builder.WriteString("\n")
		}
	}

	return builder.String(), nil
}

// GetMetadata 获取 DOCX 文件的元数据
func (r *DocxReader) GetMetadata(filePath string) (map[string]string, error) {
	zipReader, err := zip.OpenReader(filePath)
	if err != nil {
		return nil, WrapError("DocxReader.GetMetadata", filePath, ErrFileOpen)
	}
	defer zipReader.Close()

	metadata := make(map[string]string)

	// 读取核心属性
	for _, file := range zipReader.File {
		if file.Name == "docProps/core.xml" {
			rc, err := file.Open()
			if err != nil {
				continue
			}
			data, err := io.ReadAll(rc)
			rc.Close()
			if err != nil {
				continue
			}

			var props CoreProperties
			if err := xml.Unmarshal(data, &props); err == nil {
				metadata["title"] = props.Title
				metadata["subject"] = props.Subject
				metadata["creator"] = props.Creator
				metadata["description"] = props.Description
				metadata["created"] = props.Created
				metadata["modified"] = props.Modified
			}
			break
		}
	}

	return metadata, nil
}

// ReadWithConfig 根据配置读取 DOCX 文件，返回结构化结果
// DOCX 文件以段落为单位，将每个段落视为一行
func (r *DocxReader) ReadWithConfig(filePath string, config *ReadConfig) (*DocumentResult, error) {
	zipReader, err := zip.OpenReader(filePath)
	if err != nil {
		return nil, WrapError("DocxReader.ReadWithConfig", filePath, ErrFileOpen)
	}
	defer zipReader.Close()

	// 查找并读取 document.xml
	var documentXML []byte
	for _, file := range zipReader.File {
		if file.Name == "word/document.xml" {
			rc, err := file.Open()
			if err != nil {
				return nil, WrapError("DocxReader.ReadWithConfig", filePath, ErrFileRead)
			}
			documentXML, err = io.ReadAll(rc)
			rc.Close()
			if err != nil {
				return nil, WrapError("DocxReader.ReadWithConfig", filePath, ErrFileRead)
			}
			break
		}
	}

	if documentXML == nil {
		return nil, WrapError("DocxReader.ReadWithConfig", filePath, ErrInvalidFormat)
	}

	// 解析 XML
	var doc WordDocument
	if err := xml.Unmarshal(documentXML, &doc); err != nil {
		return nil, WrapError("DocxReader.ReadWithConfig", filePath, ErrFileParse)
	}

	result := &DocumentResult{
		FilePath:   filePath,
		TotalPages: 1, // DOCX 作为单页处理
		Pages:      make([]PageContent, 0),
		Metadata:   make(map[string]string),
	}

	// 获取元数据
	metadata, _ := r.GetMetadata(filePath)
	result.Metadata = metadata

	// 提取所有段落和表格行
	lines := make([]string, 0)

	// 提取段落文本
	for _, para := range doc.Body.Paragraphs {
		var lineBuilder strings.Builder
		for _, run := range para.Runs {
			lineBuilder.WriteString(run.Text)
		}
		line := lineBuilder.String()
		if line != "" {
			lines = append(lines, line)
		}
	}

	// 提取表格文本
	for _, table := range doc.Body.Tables {
		for _, row := range table.Rows {
			var rowBuilder strings.Builder
			for cellIndex, cell := range row.Cells {
				if cellIndex > 0 {
					rowBuilder.WriteString("\t")
				}
				for _, para := range cell.Paragraphs {
					for _, run := range para.Runs {
						rowBuilder.WriteString(run.Text)
						rowBuilder.WriteString(" ")
					}
				}
			}
			line := strings.TrimSpace(rowBuilder.String())
			if line != "" {
				lines = append(lines, line)
			}
		}
	}

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
