package docreader

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
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
		return "", fmt.Errorf("failed to open docx file: %w", err)
	}
	defer zipReader.Close()

	// 查找并读取 document.xml
	var documentXML []byte
	for _, file := range zipReader.File {
		if file.Name == "word/document.xml" {
			rc, err := file.Open()
			if err != nil {
				return "", fmt.Errorf("failed to open document.xml: %w", err)
			}
			documentXML, err = io.ReadAll(rc)
			rc.Close()
			if err != nil {
				return "", fmt.Errorf("failed to read document.xml: %w", err)
			}
			break
		}
	}

	if documentXML == nil {
		return "", fmt.Errorf("document.xml not found")
	}

	// 解析 XML
	var doc WordDocument
	if err := xml.Unmarshal(documentXML, &doc); err != nil {
		return "", fmt.Errorf("failed to parse XML: %w", err)
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
		return nil, fmt.Errorf("failed to open docx file: %w", err)
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
