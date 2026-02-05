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
	t.Run("不存在的文件", func(t *testing.T) {
		_, err := ReadDocument("nonexistent.docx")
		if err == nil {
			t.Fatal("期望出现错误")
		}
		if !strings.Contains(err.Error(), "file not found") {
			t.Errorf("期望错误信息包含 'file not found'，实际: %v", err)
		}
	})
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
			if _, err := reader.ReadText("nonexistent.file"); err == nil {
				t.Error("期望读取不存在文件时出现错误")
			}

			if _, err := reader.GetMetadata("nonexistent.file"); err == nil {
				t.Error("期望获取不存在文件元数据时出现错误")
			}
		})
	}
}

// testFileReader 通用的文件读取测试辅助函数
func testFileReader(t *testing.T, filename string, reader DocumentReader, extraTests func(*testing.T, string, string, map[string]string)) {
	t.Helper()

	testFile := filepath.Join("testdata", filename)
	if _, err := os.Stat(testFile); err != nil {
		t.Skipf("测试文件不存在: %s", filename)
	}

	fileInfo, _ := os.Stat(testFile)

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

	// 输出基本统计信息
	t.Logf("=== %s 文件统计 ===", strings.ToUpper(filepath.Ext(filename)[1:]))
	t.Logf("文件大小: %s", formatFileSize(fileInfo.Size()))
	t.Logf("内容长度: %d 字符", len(content))
	t.Logf("处理时间: %v", duration)

	// 执行额外的测试
	if extraTests != nil {
		extraTests(t, testFile, content, metadata)
	}

	t.Logf("元数据: %+v", metadata)
}

// formatFileSize 格式化文件大小
func formatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

// TestDocxReaderWithRealFile 测试 DOCX 读取器
func TestDocxReaderWithRealFile(t *testing.T) {
	testFileReader(t, "test.docx", &DocxReader{}, nil)
}

// TestPdfReaderWithRealFile 测试 PDF 读取器
func TestPdfReaderWithRealFile(t *testing.T) {
	testFileReader(t, "test.pdf", &PdfReader{}, func(t *testing.T, path, content string, metadata map[string]string) {
		if pages, ok := metadata["pages"]; ok {
			t.Logf("页数: %s", pages)
		}
	})
}

// TestXlsxReaderWithRealFile 测试 XLSX 读取器
func TestXlsxReaderWithRealFile(t *testing.T) {
	testFileReader(t, "test.xlsx", &XlsxReader{}, func(t *testing.T, path, content string, metadata map[string]string) {
		reader := &XlsxReader{}
		if allData, err := reader.GetAllSheetsData(path); err == nil {
			totalRows := 0
			for _, rows := range allData {
				totalRows += len(rows)
			}
			t.Logf("总行数: %d", totalRows)
		}
		if sheets, ok := metadata["sheet_count"]; ok {
			t.Logf("工作表数量: %s", sheets)
		}
	})
}

// TestPptxReaderWithRealFile 测试 PPTX 读取器
func TestPptxReaderWithRealFile(t *testing.T) {
	testFileReader(t, "test.pptx", &PptxReader{}, func(t *testing.T, path, content string, metadata map[string]string) {
		reader := &PptxReader{}
		if slides, err := reader.GetSlides(path); err == nil {
			t.Logf("幻灯片数量: %d", len(slides))
		}
	})
}

// TestTxtReaderWithRealFile 测试 TXT 读取器
func TestTxtReaderWithRealFile(t *testing.T) {
	testFileReader(t, "test.txt", &TxtReader{}, func(t *testing.T, path, content string, metadata map[string]string) {
		t.Logf("行数: %d", strings.Count(content, "\n")+1)
	})
}

// TestCsvReaderWithRealFile 测试 CSV 读取器
func TestCsvReaderWithRealFile(t *testing.T) {
	testFileReader(t, "test.csv", &CsvReader{}, func(t *testing.T, path, content string, metadata map[string]string) {
		reader := &CsvReader{}
		if records, err := reader.GetRecords(path); err == nil {
			t.Logf("记录数量: %d 行", len(records))
			if len(records) > 0 {
				t.Logf("列数: %d", len(records[0]))
			}
		}
	})
}

// TestMdReaderWithRealFile 测试 Markdown 读取器
func TestMdReaderWithRealFile(t *testing.T) {
	testFileReader(t, "test.md", &MdReader{}, func(t *testing.T, path, content string, metadata map[string]string) {
		t.Logf("行数: %d", strings.Count(content, "\n")+1)
	})
}

// TestRtfReaderWithRealFile 测试 RTF 读取器
func TestRtfReaderWithRealFile(t *testing.T) {
	testFileReader(t, "test.rtf", &RtfReader{}, nil)
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
		tmpFile := filepath.Join("testdata", "test.unknown")
		if err := os.WriteFile(tmpFile, []byte("test"), 0644); err != nil {
			t.Fatalf("创建临时文件失败: %v", err)
		}
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

	if len(formats) == 0 {
		t.Fatal("支持的格式列表不应为空")
	}

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

	// 验证返回的是副本
	originalFirst := formats[0]
	formats[0] = ".test"
	formats2 := GetSupportedFormats()
	if formats2[0] != originalFirst {
		t.Error("GetSupportedFormats 应该返回副本")
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
		{"docx", true},
		{"DOCX", true},
		{"PDF", true},
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

		t.Logf("%-8s %-12s %-15d %-15v",
			format,
			formatFileSize(fileInfo.Size()),
			len(doc.Content),
			duration)
	}
}

// BenchmarkReadDocument 性能基准测试
func BenchmarkReadDocument(b *testing.B) {
	testFiles := map[string]string{
		"DOCX": "testdata/test.docx",
		"PDF":  "testdata/test.pdf",
		"XLSX": "testdata/test.xlsx",
		"PPTX": "testdata/test.pptx",
		"TXT":  "testdata/test.txt",
		"CSV":  "testdata/test.csv",
		"MD":   "testdata/test.md",
	}

	for name, testFile := range testFiles {
		if _, err := os.Stat(testFile); err != nil {
			continue
		}

		b.Run(name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = ReadDocument(testFile)
			}
		})
	}
}
