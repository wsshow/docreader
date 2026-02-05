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

// pageLineFilter 存储单页的行过滤配置
type pageLineFilter struct {
	lines   map[int]bool // 要读取的行号集合
	readAll bool         // 是否读取所有行
}

// buildPageLineMap 构建页码到行配置的映射
func buildPageLineMap(config *ReadConfig, totalPages int) map[int]pageLineFilter {
	result := make(map[int]pageLineFilter)

	// 如果有详细的页面配置，优先使用
	if config != nil && len(config.PageConfigs) > 0 {
		for _, pageConfig := range config.PageConfigs {
			if pageConfig.PageIndex < 0 || pageConfig.PageIndex >= totalPages {
				continue
			}

			linesSet := make(map[int]bool)

			// 添加离散的行号
			for _, line := range pageConfig.LineSelector.Indexes {
				if line >= 0 {
					linesSet[line] = true
				}
			}

			// 添加行范围
			for _, lineRange := range pageConfig.LineSelector.Ranges {
				start, end := lineRange[0], lineRange[1]
				if start < 0 {
					start = 0
				}
				for i := start; i <= end; i++ {
					linesSet[i] = true
				}
			}

			result[pageConfig.PageIndex] = pageLineFilter{
				lines:   linesSet,
				readAll: len(linesSet) == 0,
			}
		}
		return result
	}

	// 使用全局配置
	// 确定要读取的页码
	pagesToRead := determinePagesToRead(config, totalPages)

	// 构建全局行配置
	var globalLineFilter pageLineFilter
	if config == nil || (config.LineSelector.Indexes == nil && config.LineSelector.Ranges == nil) {
		globalLineFilter = pageLineFilter{readAll: true}
	} else {
		linesSet := make(map[int]bool)

		// 添加离散的行号
		for _, line := range config.LineSelector.Indexes {
			if line >= 0 {
				linesSet[line] = true
			}
		}

		// 添加行号范围
		for _, lineRange := range config.LineSelector.Ranges {
			start, end := lineRange[0], lineRange[1]
			if start < 0 {
				start = 0
			}
			for i := start; i <= end; i++ {
				linesSet[i] = true
			}
		}

		globalLineFilter = pageLineFilter{
			lines:   linesSet,
			readAll: len(linesSet) == 0,
		}
	}

	// 将全局配置应用到所有要读取的页
	for _, pageIndex := range pagesToRead {
		result[pageIndex] = globalLineFilter
	}

	return result
}

// filterLinesForPage 根据页面配置筛选行
func filterLinesForPage(lines []string, filter pageLineFilter) []string {
	if filter.readAll {
		return lines
	}

	result := make([]string, 0)
	for i := 0; i < len(lines); i++ {
		if filter.lines[i] {
			result = append(result, lines[i])
		}
	}

	return result
}

// filterLinesForSinglePage 为单页文档筛选行（用于 TXT/MD/CSV/RTF/DOCX）
func filterLinesForSinglePage(lines []string, config *ReadConfig) []string {
	if config != nil && len(config.PageConfigs) > 0 {
		// 查找页面0的配置
		for _, pageConfig := range config.PageConfigs {
			if pageConfig.PageIndex == 0 {
				linesSet := make(map[int]bool)

				// 添加离散行号
				for _, line := range pageConfig.LineSelector.Indexes {
					if line >= 0 {
						linesSet[line] = true
					}
				}

				// 添加行范围
				for _, lineRange := range pageConfig.LineSelector.Ranges {
					start, end := lineRange[0], lineRange[1]
					if start < 0 {
						start = 0
					}
					for i := start; i <= end; i++ {
						linesSet[i] = true
					}
				}

				filter := pageLineFilter{
					lines:   linesSet,
					readAll: len(linesSet) == 0,
				}

				return filterLinesForPage(lines, filter)
			}
		}
		return []string{}
	}

	// 使用全局配置
	pageLineMap := buildPageLineMap(config, 1)
	if filter, ok := pageLineMap[0]; ok {
		return filterLinesForPage(lines, filter)
	}
	return lines
}

// determinePagesToRead 根据配置确定要读取的页码（索引从0开始）
func determinePagesToRead(config *ReadConfig, totalPages int) []int {
	if config == nil {
		// 如果没有配置，读取所有页
		pages := make([]int, totalPages)
		for i := 0; i < totalPages; i++ {
			pages[i] = i
		}
		return pages
	}

	pagesSet := make(map[int]bool)

	// 添加离散页码
	for _, p := range config.PageSelector.Indexes {
		if p >= 0 && p < totalPages {
			pagesSet[p] = true
		}
	}

	// 添加页码范围
	for _, pageRange := range config.PageSelector.Ranges {
		start, end := pageRange[0], pageRange[1]
		if start < 0 {
			start = 0
		}
		if end >= totalPages {
			end = totalPages - 1
		}
		for i := start; i <= end; i++ {
			pagesSet[i] = true
		}
	}

	// 如果没有指定任何页码，返回所有页
	if len(pagesSet) == 0 {
		pages := make([]int, totalPages)
		for i := 0; i < totalPages; i++ {
			pages[i] = i
		}
		return pages
	}

	// 转换为有序切片
	pages := make([]int, 0, len(pagesSet))
	for i := 0; i < totalPages; i++ {
		if pagesSet[i] {
			pages = append(pages, i)
		}
	}

	return pages
}
