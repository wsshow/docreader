package docreader

import (
	"os"
	"path/filepath"
	"testing"
)

// TestReadDocument 测试统一文档读取接口
func TestReadDocument(t *testing.T) {
	tests := []struct {
		name        string
		filepath    string
		shouldError bool
	}{
		{
			name:        "不支持的格式",
			filepath:    "test.txt",
			shouldError: true,
		},
		{
			name:        "不存在的文件",
			filepath:    "nonexistent.docx",
			shouldError: true,
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
		})
	}
}

// TestDocxReader 测试 DOCX 读取器
func TestDocxReader(t *testing.T) {
	reader := &DocxReader{}

	t.Run("读取不存在的文件", func(t *testing.T) {
		_, err := reader.ReadText("nonexistent.docx")
		if err == nil {
			t.Error("期望出现错误")
		}
	})

	t.Run("获取不存在文件的元数据", func(t *testing.T) {
		_, err := reader.GetMetadata("nonexistent.docx")
		if err == nil {
			t.Error("期望出现错误")
		}
	})
}

// TestPdfReader 测试 PDF 读取器
func TestPdfReader(t *testing.T) {
	reader := &PdfReader{}

	t.Run("读取不存在的文件", func(t *testing.T) {
		_, err := reader.ReadText("nonexistent.pdf")
		if err == nil {
			t.Error("期望出现错误")
		}
	})

	t.Run("获取不存在文件的元数据", func(t *testing.T) {
		_, err := reader.GetMetadata("nonexistent.pdf")
		if err == nil {
			t.Error("期望出现错误")
		}
	})
}

// TestXlsxReader 测试 XLSX 读取器
func TestXlsxReader(t *testing.T) {
	reader := &XlsxReader{}

	t.Run("读取不存在的文件", func(t *testing.T) {
		_, err := reader.ReadText("nonexistent.xlsx")
		if err == nil {
			t.Error("期望出现错误")
		}
	})

	t.Run("获取不存在文件的元数据", func(t *testing.T) {
		_, err := reader.GetMetadata("nonexistent.xlsx")
		if err == nil {
			t.Error("期望出现错误")
		}
	})

	t.Run("获取不存在文件的工作表数据", func(t *testing.T) {
		_, err := reader.GetSheetData("nonexistent.xlsx", "Sheet1")
		if err == nil {
			t.Error("期望出现错误")
		}
	})
}

// TestPptxReader 测试 PPTX 读取器
func TestPptxReader(t *testing.T) {
	reader := &PptxReader{}

	t.Run("读取不存在的文件", func(t *testing.T) {
		_, err := reader.ReadText("nonexistent.pptx")
		if err == nil {
			t.Error("期望出现错误")
		}
	})

	t.Run("获取不存在文件的元数据", func(t *testing.T) {
		_, err := reader.GetMetadata("nonexistent.pptx")
		if err == nil {
			t.Error("期望出现错误")
		}
	})

	t.Run("获取不存在文件的幻灯片", func(t *testing.T) {
		_, err := reader.GetSlides("nonexistent.pptx")
		if err == nil {
			t.Error("期望出现错误")
		}
	})
}

// 集成测试
func TestIntegrationWithRealFiles(t *testing.T) {
	testDir := "testdata"

	// 测试 DOCX
	docxPath := filepath.Join(testDir, "test.docx")
	if _, err := os.Stat(docxPath); err == nil {
		doc, err := ReadDocument(docxPath)
		if err != nil {
			t.Errorf("读取 DOCX 失败: %v", err)
		}
		if doc.Content == "" {
			t.Error("DOCX 内容为空")
		}
		t.Logf("DOCX 内容长度: %d, 内容: %s", len(doc.Content), doc.Content)
	}

	// 测试 PDF
	pdfPath := filepath.Join(testDir, "test.pdf")
	if _, err := os.Stat(pdfPath); err == nil {
		doc, err := ReadDocument(pdfPath)
		if err != nil {
			t.Errorf("读取 PDF 失败: %v", err)
		}
		if doc.Content == "" {
			t.Error("PDF 内容为空")
		}
		t.Logf("PDF 内容长度: %d, 内容: %s", len(doc.Content), doc.Content)
	}

	// 测试 XLSX
	xlsxPath := filepath.Join(testDir, "test.xlsx")
	if _, err := os.Stat(xlsxPath); err == nil {
		doc, err := ReadDocument(xlsxPath)
		if err != nil {
			t.Errorf("读取 XLSX 失败: %v", err)
		}
		if doc.Content == "" {
			t.Error("XLSX 内容为空")
		}
		t.Logf("XLSX 内容长度: %d, 内容: %s", len(doc.Content), doc.Content)
	}

	// 测试 PPTX
	pptxPath := filepath.Join(testDir, "test.pptx")
	if _, err := os.Stat(pptxPath); err == nil {
		doc, err := ReadDocument(pptxPath)
		if err != nil {
			t.Errorf("读取 PPTX 失败: %v", err)
		}
		if doc.Content == "" {
			t.Error("PPTX 内容为空")
		}
		t.Logf("PPTX 内容长度: %d, 内容: %s", len(doc.Content), doc.Content)
	}
}

// BenchmarkReadDocument 性能基准测试
func BenchmarkReadDocument(b *testing.B) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		b.Skip("跳过性能测试")
	}

	testFile := "testdata/test.docx"
	if _, err := os.Stat(testFile); err != nil {
		b.Skip("测试文件不存在")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ReadDocument(testFile)
	}
}
