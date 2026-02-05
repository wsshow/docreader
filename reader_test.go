package docreader

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestReadDocument 测试统一文档读取接口
func TestReadDocument(t *testing.T) {
	tests := []struct {
		name        string
		filepath    string
		shouldError bool
		errorMsg    string
	}{
		{
			name:        "不存在的文件",
			filepath:    "nonexistent.docx",
			shouldError: true,
			errorMsg:    "file not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ReadDocument(tt.filepath)
			if tt.shouldError && err == nil {
				t.Error("期望出现错误，但没有错误")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("不期望出现错误，但得到: %v", err)
			}
			if tt.shouldError && err != nil && tt.errorMsg != "" {
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("期望错误信息包含 '%s'，但得到: %v", tt.errorMsg, err)
				}
			}
		})
	}
}

// TestErrorHandling 统一测试所有读取器的错误处理
func TestErrorHandling(t *testing.T) {
	readers := map[string]DocumentReader{
		"DOCX": &DocxReader{},
		"PDF":  &PdfReader{},
		"XLSX": &XlsxReader{},
		"PPTX": &PptxReader{},
		"TXT":  &TxtReader{},
		"CSV":  &CsvReader{},
		"MD":   &MdReader{},
		"RTF":  &RtfReader{},
	}

	for name, reader := range readers {
		t.Run(name, func(t *testing.T) {
			_, err := reader.ReadText("nonexistent.file")
			if err == nil {
				t.Errorf("%s: 期望读取不存在文件时出现错误", name)
			}

			_, err = reader.GetMetadata("nonexistent.file")
			if err == nil {
				t.Errorf("%s: 期望获取不存在文件元数据时出现错误", name)
			}
		})
	}
}

// TestDocxReaderWithRealFile 测试 DOCX 读取器
func TestDocxReaderWithRealFile(t *testing.T) {
	testFile := filepath.Join("testdata", "test.docx")
	if _, err := os.Stat(testFile); err != nil {
		t.Skip("测试文件不存在: test.docx")
	}

	fileInfo, _ := os.Stat(testFile)
	reader := &DocxReader{}

	// 测试读取文本
	start := time.Now()
	content, err := reader.ReadText(testFile)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("读取失败: %v", err)
	}

	// 获取元数据
	metadata, err := reader.GetMetadata(testFile)
	if err != nil {
		t.Errorf("获取元数据失败: %v", err)
	}

	// 输出统计信息
	t.Logf("=== DOCX 文件统计 ===")
	t.Logf("文件大小: %d 字节 (%.2f KB)", fileInfo.Size(), float64(fileInfo.Size())/1024)
	t.Logf("内容长度: %d 字符", len(content))
	t.Logf("处理时间: %v", duration)
	t.Logf("元数据: %+v", metadata)
}

// TestPdfReaderWithRealFile 测试 PDF 读取器
func TestPdfReaderWithRealFile(t *testing.T) {
	testFile := filepath.Join("testdata", "test.pdf")
	if _, err := os.Stat(testFile); err != nil {
		t.Skip("测试文件不存在: test.pdf")
	}

	fileInfo, _ := os.Stat(testFile)
	reader := &PdfReader{}

	// 测试读取文本
	start := time.Now()
	content, err := reader.ReadText(testFile)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("读取失败: %v", err)
	}

	// 获取元数据
	metadata, err := reader.GetMetadata(testFile)
	if err != nil {
		t.Errorf("获取元数据失败: %v", err)
	}

	// 输出统计信息
	t.Logf("=== PDF 文件统计 ===")
	t.Logf("文件大小: %d 字节 (%.2f MB)", fileInfo.Size(), float64(fileInfo.Size())/1024/1024)
	t.Logf("内容长度: %d 字符", len(content))
	t.Logf("处理时间: %v", duration)
	if pages, ok := metadata["pages"]; ok {
		t.Logf("页数: %s", pages)
	}
	t.Logf("元数据: %+v", metadata)
}

// TestXlsxReaderWithRealFile 测试 XLSX 读取器
func TestXlsxReaderWithRealFile(t *testing.T) {
	testFile := filepath.Join("testdata", "test.xlsx")
	if _, err := os.Stat(testFile); err != nil {
		t.Skip("测试文件不存在: test.xlsx")
	}

	fileInfo, _ := os.Stat(testFile)
	reader := &XlsxReader{}

	// 测试读取文本
	start := time.Now()
	content, err := reader.ReadText(testFile)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("读取失败: %v", err)
	}

	// 获取元数据
	metadata, err := reader.GetMetadata(testFile)
	if err != nil {
		t.Errorf("获取元数据失败: %v", err)
	}

	// 获取工作表数据
	allData, err := reader.GetAllSheetsData(testFile)
	totalRows := 0
	if err == nil {
		for _, rows := range allData {
			totalRows += len(rows)
		}
	}

	// 输出统计信息
	t.Logf("=== XLSX 文件统计 ===")
	t.Logf("文件大小: %d 字节 (%.2f KB)", fileInfo.Size(), float64(fileInfo.Size())/1024)
	t.Logf("内容长度: %d 字符", len(content))
	t.Logf("处理时间: %v", duration)
	if sheets, ok := metadata["sheet_count"]; ok {
		t.Logf("工作表数量: %s", sheets)
	}
	t.Logf("总行数: %d", totalRows)
	t.Logf("元数据: %+v", metadata)
}

// TestPptxReaderWithRealFile 测试 PPTX 读取器
func TestPptxReaderWithRealFile(t *testing.T) {
	testFile := filepath.Join("testdata", "test.pptx")
	if _, err := os.Stat(testFile); err != nil {
		t.Skip("测试文件不存在: test.pptx")
	}

	fileInfo, _ := os.Stat(testFile)
	reader := &PptxReader{}

	// 测试读取文本
	start := time.Now()
	content, err := reader.ReadText(testFile)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("读取失败: %v", err)
	}

	// 获取元数据
	metadata, err := reader.GetMetadata(testFile)
	if err != nil {
		t.Errorf("获取元数据失败: %v", err)
	}

	// 获取幻灯片
	slides, err := reader.GetSlides(testFile)
	if err != nil {
		t.Errorf("获取幻灯片失败: %v", err)
	}

	// 输出统计信息
	t.Logf("=== PPTX 文件统计 ===")
	t.Logf("文件大小: %d 字节 (%.2f KB)", fileInfo.Size(), float64(fileInfo.Size())/1024)
	t.Logf("内容长度: %d 字符", len(content))
	t.Logf("处理时间: %v", duration)
	t.Logf("幻灯片数量: %d", len(slides))
	t.Logf("元数据: %+v", metadata)
}

// TestTxtReaderWithRealFile 测试 TXT 读取器
func TestTxtReaderWithRealFile(t *testing.T) {
	testFile := filepath.Join("testdata", "test.txt")
	if _, err := os.Stat(testFile); err != nil {
		t.Skip("测试文件不存在: test.txt")
	}

	fileInfo, _ := os.Stat(testFile)
	reader := &TxtReader{}

	// 测试读取文本
	start := time.Now()
	content, err := reader.ReadText(testFile)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("读取失败: %v", err)
	}

	// 获取元数据
	metadata, err := reader.GetMetadata(testFile)
	if err != nil {
		t.Errorf("获取元数据失败: %v", err)
	}

	// 输出统计信息
	t.Logf("=== TXT 文件统计 ===")
	t.Logf("文件大小: %d 字节", fileInfo.Size())
	t.Logf("内容长度: %d 字符", len(content))
	t.Logf("行数: %d", strings.Count(content, "\n")+1)
	t.Logf("处理时间: %v", duration)
	t.Logf("元数据: %+v", metadata)
}

// TestCsvReaderWithRealFile 测试 CSV 读取器
func TestCsvReaderWithRealFile(t *testing.T) {
	testFile := filepath.Join("testdata", "test.csv")
	if _, err := os.Stat(testFile); err != nil {
		t.Skip("测试文件不存在: test.csv")
	}

	fileInfo, _ := os.Stat(testFile)
	reader := &CsvReader{}

	// 测试读取文本
	start := time.Now()
	content, err := reader.ReadText(testFile)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("读取失败: %v", err)
	}

	// 获取元数据
	metadata, err := reader.GetMetadata(testFile)
	if err != nil {
		t.Errorf("获取元数据失败: %v", err)
	}

	// 获取记录
	records, err := reader.GetRecords(testFile)
	if err != nil {
		t.Errorf("获取记录失败: %v", err)
	}

	// 输出统计信息
	t.Logf("=== CSV 文件统计 ===")
	t.Logf("文件大小: %d 字节", fileInfo.Size())
	t.Logf("内容长度: %d 字符", len(content))
	t.Logf("处理时间: %v", duration)
	t.Logf("记录数量: %d 行", len(records))
	if len(records) > 0 {
		t.Logf("列数: %d", len(records[0]))
	}
	t.Logf("元数据: %+v", metadata)
}

// TestMdReaderWithRealFile 测试 Markdown 读取器
func TestMdReaderWithRealFile(t *testing.T) {
	testFile := filepath.Join("testdata", "test.md")
	if _, err := os.Stat(testFile); err != nil {
		t.Skip("测试文件不存在: test.md")
	}

	fileInfo, _ := os.Stat(testFile)
	reader := &MdReader{}

	// 测试读取文本
	start := time.Now()
	content, err := reader.ReadText(testFile)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("读取失败: %v", err)
	}

	// 获取元数据
	metadata, err := reader.GetMetadata(testFile)
	if err != nil {
		t.Errorf("获取元数据失败: %v", err)
	}

	// 输出统计信息
	t.Logf("=== Markdown 文件统计 ===")
	t.Logf("文件大小: %d 字节", fileInfo.Size())
	t.Logf("内容长度: %d 字符", len(content))
	t.Logf("行数: %d", strings.Count(content, "\n")+1)
	t.Logf("处理时间: %v", duration)
	t.Logf("元数据: %+v", metadata)
}

// TestRtfReaderWithRealFile 测试 RTF 读取器
func TestRtfReaderWithRealFile(t *testing.T) {
	testFile := filepath.Join("testdata", "test.rtf")
	if _, err := os.Stat(testFile); err != nil {
		t.Skip("测试文件不存在: test.rtf")
	}

	fileInfo, _ := os.Stat(testFile)
	reader := &RtfReader{}

	// 测试读取文本
	start := time.Now()
	content, err := reader.ReadText(testFile)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("读取失败: %v", err)
	}

	// 获取元数据
	metadata, err := reader.GetMetadata(testFile)
	if err != nil {
		t.Errorf("获取元数据失败: %v", err)
	}

	// 输出统计信息
	t.Logf("=== RTF 文件统计 ===")
	t.Logf("文件大小: %d 字节", fileInfo.Size())
	t.Logf("提取内容长度: %d 字符", len(content))
	t.Logf("处理时间: %v", duration)
	t.Logf("元数据: %+v", metadata)
}

// TestFormatDetection 测试格式检测
func TestFormatDetection(t *testing.T) {
	tests := []struct {
		filename string
		expected string
	}{
		{"test.docx", ".docx"},
		{"test.DOCX", ".docx"},
		{"test.pdf", ".pdf"},
		{"test.PDF", ".pdf"},
		{"test.xlsx", ".xlsx"},
		{"test.pptx", ".pptx"},
		{"test.txt", ".txt"},
		{"test.csv", ".csv"},
		{"test.md", ".md"},
		{"test.markdown", ".markdown"},
		{"test.rtf", ".rtf"},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			ext := strings.ToLower(filepath.Ext(tt.filename))
			if ext != tt.expected {
				t.Errorf("期望扩展名 %s，得到 %s", tt.expected, ext)
			}
		})
	}
}

// TestErrorTypes 测试错误类型
func TestErrorTypes(t *testing.T) {
	t.Run("文件不存在错误", func(t *testing.T) {
		_, err := ReadDocument("nonexistent.docx")
		if err == nil {
			t.Fatal("期望出现错误")
		}
		if !IsFileNotFound(err) {
			t.Errorf("期望 FileNotFound 错误，得到: %v", err)
		}
	})

	t.Run("不支持的格式错误", func(t *testing.T) {
		// 创建一个临时文件
		tmpFile := filepath.Join("testdata", "test.unknown")
		os.WriteFile(tmpFile, []byte("test"), 0644)
		defer os.Remove(tmpFile)

		_, err := ReadDocument(tmpFile)
		if err == nil {
			t.Fatal("期望出现错误")
		}
		if !IsUnsupportedFormat(err) {
			t.Errorf("期望 UnsupportedFormat 错误，得到: %v", err)
		}
	})
}

// TestGetSupportedFormats 测试获取支持的格式列表
func TestGetSupportedFormats(t *testing.T) {
	formats := GetSupportedFormats()

	// 检查返回的格式列表不为空
	if len(formats) == 0 {
		t.Error("支持的格式列表不应为空")
	}

	// 检查是否包含常见格式
	expectedFormats := []string{".docx", ".pdf", ".xlsx", ".pptx", ".txt", ".csv", ".md", ".rtf"}
	for _, expected := range expectedFormats {
		found := false
		for _, format := range formats {
			if format == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("期望格式列表包含 %s", expected)
		}
	}

	t.Logf("支持的格式: %v", formats)

	// 验证返回的是副本（修改不影响原始数据）
	originalFirst := formats[0]
	formats[0] = ".test"
	formats2 := GetSupportedFormats()
	if formats2[0] != originalFirst {
		t.Error("GetSupportedFormats 应该返回副本，而不是原始切片")
	}
}

// TestIsFormatSupported 测试格式支持检查
func TestIsFormatSupported(t *testing.T) {
	tests := []struct {
		format   string
		expected bool
	}{
		{".docx", true},
		{".pdf", true},
		{".xlsx", true},
		{".pptx", true},
		{".txt", true},
		{".csv", true},
		{".md", true},
		{".markdown", true},
		{".rtf", true},
		{".doc", false},
		{".xls", false},
		{".ppt", false},
		{".unknown", false},
		{"docx", true}, // 不带点号
		{"DOCX", true}, // 大写
		{"PDF", true},  // 大写
	}

	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			result := IsFormatSupported(tt.format)
			if result != tt.expected {
				t.Errorf("IsFormatSupported(%s) = %v, 期望 %v", tt.format, result, tt.expected)
			}
		})
	}
}

// TestAllFormatsPerformance 测试所有格式的性能对比
func TestAllFormatsPerformance(t *testing.T) {
	testFiles := map[string]string{
		"DOCX": "testdata/test.docx",
		"PDF":  "testdata/test.pdf",
		"XLSX": "testdata/test.xlsx",
		"PPTX": "testdata/test.pptx",
		"TXT":  "testdata/test.txt",
		"CSV":  "testdata/test.csv",
		"MD":   "testdata/test.md",
	}

	t.Log("=== 所有格式性能对比 ===")
	t.Logf("%-8s %-12s %-15s %-15s", "格式", "文件大小", "内容长度", "处理时间")
	t.Log(strings.Repeat("-", 60))

	for format, filePath := range testFiles {
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			continue
		}

		start := time.Now()
		doc, err := ReadDocument(filePath)
		duration := time.Since(start)

		if err != nil {
			t.Logf("%-8s 读取失败: %v", format, err)
			continue
		}

		sizeStr := fmt.Sprintf("%.2f KB", float64(fileInfo.Size())/1024)
		if fileInfo.Size() > 1024*1024 {
			sizeStr = fmt.Sprintf("%.2f MB", float64(fileInfo.Size())/1024/1024)
		}

		t.Logf("%-8s %-12s %-15d %-15v",
			format,
			sizeStr,
			len(doc.Content),
			duration)
	}
}

// BenchmarkReadDocument 性能基准测试
func BenchmarkReadDocument(b *testing.B) {
	testFiles := []string{
		"testdata/test.docx",
		"testdata/test.pdf",
		"testdata/test.xlsx",
		"testdata/test.pptx",
		"testdata/test.txt",
		"testdata/test.csv",
		"testdata/test.md",
	}

	for _, testFile := range testFiles {
		if _, err := os.Stat(testFile); err != nil {
			continue
		}

		b.Run(filepath.Base(testFile), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = ReadDocument(testFile)
			}
		})
	}
}
