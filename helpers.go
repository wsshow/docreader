package docreader

// helpers.go 包含文档读取的公共辅助函数
// 这些函数被多个格式读取器共享使用

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
	globalLineFilter := buildGlobalLineFilter(config)

	// 将全局配置应用到所有要读取的页
	for _, pageIndex := range pagesToRead {
		result[pageIndex] = globalLineFilter
	}

	return result
}

// buildGlobalLineFilter 构建全局行过滤器
func buildGlobalLineFilter(config *ReadConfig) pageLineFilter {
	if config == nil || (config.LineSelector.Indexes == nil && config.LineSelector.Ranges == nil) {
		return pageLineFilter{readAll: true}
	}

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

	return pageLineFilter{
		lines:   linesSet,
		readAll: len(linesSet) == 0,
	}
}

// filterLinesForPage 根据页面配置筛选行
func filterLinesForPage(lines []string, filter pageLineFilter) []string {
	if filter.readAll {
		return lines
	}

	result := make([]string, 0, len(filter.lines))
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
		return makeAllPagesSlice(totalPages)
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
		return makeAllPagesSlice(totalPages)
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

// makeAllPagesSlice 创建包含所有页码的切片
func makeAllPagesSlice(totalPages int) []int {
	pages := make([]int, totalPages)
	for i := 0; i < totalPages; i++ {
		pages[i] = i
	}
	return pages
}
