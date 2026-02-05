package docreader

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

// PptxReader 用于读取 .pptx 文件
type PptxReader struct{}

// Slide 表示幻灯片的 XML 结构
type Slide struct {
	XMLName   xml.Name `xml:"sld"`
	CommonSld struct {
		ShapeTree struct {
			Shapes []struct {
				TextBody struct {
					Paragraphs []struct {
						Runs []struct {
							Text string `xml:"t"`
						} `xml:"r"`
					} `xml:"p"`
				} `xml:"txBody"`
			} `xml:"sp"`
		} `xml:"spTree"`
	} `xml:"cSld"`
}

// PresentationProps 表示演示文稿属性
type PresentationProps struct {
	XMLName  xml.Name `xml:"coreProperties"`
	Title    string   `xml:"title"`
	Subject  string   `xml:"subject"`
	Creator  string   `xml:"creator"`
	Keywords string   `xml:"keywords"`
	Created  string   `xml:"created"`
	Modified string   `xml:"modified"`
}

// ReadText 读取 PPTX 文件的文本内容
func (r *PptxReader) ReadText(filePath string) (string, error) {
	// 打开 zip 文件
	zipReader, err := zip.OpenReader(filePath)
	if err != nil {
		return "", WrapError("PptxReader.ReadText", filePath, ErrFileOpen)
	}
	defer zipReader.Close()

	var builder strings.Builder
	slideNum := 1

	// 遍历所有文件，查找幻灯片
	for _, file := range zipReader.File {
		// 检查是否是幻灯片文件
		if strings.HasPrefix(file.Name, "ppt/slides/slide") && strings.HasSuffix(file.Name, ".xml") {
			// 读取幻灯片内容
			rc, err := file.Open()
			if err != nil {
				continue
			}

			slideXML, err := io.ReadAll(rc)
			rc.Close()
			if err != nil {
				continue
			}

			// 解析 XML
			var slide Slide
			if err := xml.Unmarshal(slideXML, &slide); err != nil {
				continue
			}

			// 提取文本
			builder.WriteString(fmt.Sprintf("\n=== 幻灯片 %d ===\n\n", slideNum))

			for _, shape := range slide.CommonSld.ShapeTree.Shapes {
				for _, para := range shape.TextBody.Paragraphs {
					for _, run := range para.Runs {
						builder.WriteString(run.Text)
					}
					builder.WriteString("\n")
				}
			}

			slideNum++
		}
	}

	if slideNum == 1 {
		return "", WrapError("PptxReader.ReadText", filePath, ErrEmptyFile)
	}

	return builder.String(), nil
}

// GetMetadata 获取 PPTX 文件的元数据
func (r *PptxReader) GetMetadata(filePath string) (map[string]string, error) {
	zipReader, err := zip.OpenReader(filePath)
	if err != nil {
		return nil, WrapError("PptxReader.GetMetadata", filePath, ErrFileOpen)
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

			var props PresentationProps
			if err := xml.Unmarshal(data, &props); err == nil {
				metadata["title"] = props.Title
				metadata["subject"] = props.Subject
				metadata["creator"] = props.Creator
				metadata["keywords"] = props.Keywords
				metadata["created"] = props.Created
				metadata["modified"] = props.Modified
			}
			break
		}
	}

	// 统计幻灯片数量
	slideCount := 0
	for _, file := range zipReader.File {
		if matched, _ := filepath.Match("ppt/slides/slide*.xml", file.Name); matched {
			slideCount++
		}
	}
	metadata["slide_count"] = fmt.Sprintf("%d", slideCount)

	return metadata, nil
}

// GetSlides 获取所有幻灯片的文本内容（按幻灯片分组）
func (r *PptxReader) GetSlides(filePath string) ([]string, error) {
	zipReader, err := zip.OpenReader(filePath)
	if err != nil {
		return nil, WrapError("PptxReader.GetSlides", filePath, ErrFileOpen)
	}
	defer zipReader.Close()

	var slides []string

	for _, file := range zipReader.File {
		if strings.HasPrefix(file.Name, "ppt/slides/slide") && strings.HasSuffix(file.Name, ".xml") {
			rc, err := file.Open()
			if err != nil {
				continue
			}

			slideXML, err := io.ReadAll(rc)
			rc.Close()
			if err != nil {
				continue
			}

			var slide Slide
			if err := xml.Unmarshal(slideXML, &slide); err != nil {
				continue
			}

			var builder strings.Builder
			for _, shape := range slide.CommonSld.ShapeTree.Shapes {
				for _, para := range shape.TextBody.Paragraphs {
					for _, run := range para.Runs {
						builder.WriteString(run.Text)
					}
					builder.WriteString("\n")
				}
			}

			slides = append(slides, builder.String())
		}
	}

	return slides, nil
}

// ReadWithConfig 根据配置读取 PPTX 文件，返回结构化结果
func (r *PptxReader) ReadWithConfig(filePath string, config *ReadConfig) (*DocumentResult, error) {
	zipReader, err := zip.OpenReader(filePath)
	if err != nil {
		return nil, WrapError("PptxReader.ReadWithConfig", filePath, ErrFileOpen)
	}
	defer zipReader.Close()

	// 先获取所有幻灯片
	type slideData struct {
		index   int
		content string
		lines   []string
	}

	allSlides := make([]slideData, 0)

	for _, file := range zipReader.File {
		if strings.HasPrefix(file.Name, "ppt/slides/slide") && strings.HasSuffix(file.Name, ".xml") {
			rc, err := file.Open()
			if err != nil {
				continue
			}

			slideXML, err := io.ReadAll(rc)
			rc.Close()
			if err != nil {
				continue
			}

			var slide Slide
			if err := xml.Unmarshal(slideXML, &slide); err != nil {
				continue
			}

			lines := make([]string, 0)
			for _, shape := range slide.CommonSld.ShapeTree.Shapes {
				for _, para := range shape.TextBody.Paragraphs {
					var lineBuilder strings.Builder
					for _, run := range para.Runs {
						lineBuilder.WriteString(run.Text)
					}
					line := lineBuilder.String()
					if line != "" {
						lines = append(lines, line)
					}
				}
			}

			allSlides = append(allSlides, slideData{
				index:   len(allSlides),
				lines:   lines,
				content: strings.Join(lines, "\n"),
			})
		}
	}

	totalSlides := len(allSlides)

	result := &DocumentResult{
		FilePath:   filePath,
		TotalPages: totalSlides,
		Pages:      make([]PageContent, 0),
		Metadata:   make(map[string]string),
	}

	// 获取元数据
	metadata, _ := r.GetMetadata(filePath)
	result.Metadata = metadata

	// 确定要读取的幻灯片和每页的行配置
	pageLineMap := buildPageLineMap(config, totalSlides)

	var contentBuilder strings.Builder
	totalLines := 0

	for slideIndex := 0; slideIndex < totalSlides; slideIndex++ {
		lineConfig, shouldRead := pageLineMap[slideIndex]
		if !shouldRead {
			continue
		}

		slide := allSlides[slideIndex]

		// 根据该页的配置筛选行
		filteredLines := filterLinesForPage(slide.lines, lineConfig)

		pageContent := PageContent{
			PageNumber: slideIndex,
			Lines:      filteredLines,
			TotalLines: len(filteredLines),
		}

		result.Pages = append(result.Pages, pageContent)
		totalLines += len(filteredLines)

		// 构建完整内容
		contentBuilder.WriteString(fmt.Sprintf("\n=== 幻灯片 %d ===\n\n", slideIndex))
		for _, line := range filteredLines {
			contentBuilder.WriteString(line)
			contentBuilder.WriteString("\n")
		}
	}

	result.TotalLines = totalLines
	result.Content = contentBuilder.String()

	return result, nil
}
