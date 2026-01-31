package docreader

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
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
			name:        "不支持的格式",
			filepath:    "test.unknown",
			shouldError: true,
			errorMsg:    "unsupported file format",
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
			if tt.shouldError && err != nil && tt.errorMsg != "" {
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("期望错误信息包含 '%s'，但得到: %v", tt.errorMsg, err)
				}
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

// TestTxtReader 测试 TXT 读取器
func TestTxtReader(t *testing.T) {
	reader := &TxtReader{}

	t.Run("读取不存在的文件", func(t *testing.T) {
		_, err := reader.ReadText("nonexistent.txt")
		if err == nil {
			t.Error("期望出现错误")
		}
	})

	t.Run("获取不存在文件的元数据", func(t *testing.T) {
		_, err := reader.GetMetadata("nonexistent.txt")
		if err == nil {
			t.Error("期望出现错误")
		}
	})
}

// TestCsvReader 测试 CSV 读取器
func TestCsvReader(t *testing.T) {
	reader := &CsvReader{}

	t.Run("读取不存在的文件", func(t *testing.T) {
		_, err := reader.ReadText("nonexistent.csv")
		if err == nil {
			t.Error("期望出现错误")
		}
	})

	t.Run("获取不存在文件的元数据", func(t *testing.T) {
		_, err := reader.GetMetadata("nonexistent.csv")
		if err == nil {
			t.Error("期望出现错误")
		}
	})

	t.Run("获取不存在文件的记录", func(t *testing.T) {
		_, err := reader.GetRecords("nonexistent.csv")
		if err == nil {
			t.Error("期望出现错误")
		}
	})
}

// TestMdReader 测试 Markdown 读取器
func TestMdReader(t *testing.T) {
	reader := &MdReader{}

	t.Run("读取不存在的文件", func(t *testing.T) {
		_, err := reader.ReadText("nonexistent.md")
		if err == nil {
			t.Error("期望出现错误")
		}
	})

	t.Run("获取不存在文件的元数据", func(t *testing.T) {
		_, err := reader.GetMetadata("nonexistent.md")
		if err == nil {
			t.Error("期望出现错误")
		}
	})
}

// TestRtfReader 测试 RTF 读取器
func TestRtfReader(t *testing.T) {
	reader := &RtfReader{}

	t.Run("读取不存在的文件", func(t *testing.T) {
		_, err := reader.ReadText("nonexistent.rtf")
		if err == nil {
			t.Error("期望出现错误")
		}
	})

	t.Run("获取不存在文件的元数据", func(t *testing.T) {
		_, err := reader.GetMetadata("nonexistent.rtf")
		if err == nil {
			t.Error("期望出现错误")
		}
	})
}

// 集成测试 - 测试所有支持的格式
func TestIntegrationWithRealFiles(t *testing.T) {
	testDir := "testdata"

	// 测试 DOCX
	t.Run("DOCX文件读取", func(t *testing.T) {
		docxPath := filepath.Join(testDir, "test.docx")
		if _, err := os.Stat(docxPath); err == nil {
			doc, err := ReadDocument(docxPath)
			if err != nil {
				t.Errorf("读取 DOCX 失败: %v", err)
				return
			}
			if doc.Content == "" {
				t.Error("DOCX 内容为空")
			}
			t.Logf("DOCX 内容长度: %d 字符", len(doc.Content))
		} else {
			t.Skip("测试文件不存在: test.docx")
		}
	})

	// 测试 PDF
	t.Run("PDF文件读取", func(t *testing.T) {
		pdfPath := filepath.Join(testDir, "test.pdf")
		if _, err := os.Stat(pdfPath); err == nil {
			doc, err := ReadDocument(pdfPath)
			if err != nil {
				t.Errorf("读取 PDF 失败: %v", err)
				return
			}
			if doc.Content == "" {
				t.Error("PDF 内容为空")
			}
			t.Logf("PDF 内容长度: %d 字符", len(doc.Content))
		} else {
			t.Skip("测试文件不存在: test.pdf")
		}
	})

	// 测试 XLSX
	t.Run("XLSX文件读取", func(t *testing.T) {
		xlsxPath := filepath.Join(testDir, "test.xlsx")
		if _, err := os.Stat(xlsxPath); err == nil {
			doc, err := ReadDocument(xlsxPath)
			if err != nil {
				t.Errorf("读取 XLSX 失败: %v", err)
				return
			}
			if doc.Content == "" {
				t.Error("XLSX 内容为空")
			}

			// 测试获取工作表数据
			reader := &XlsxReader{}
			allData, err := reader.GetAllSheetsData(xlsxPath)
			if err != nil {
				t.Errorf("获取工作表数据失败: %v", err)
			} else {
				t.Logf("XLSX 工作表数量: %d", len(allData))
			}
		} else {
			t.Skip("测试文件不存在: test.xlsx")
		}
	})

	// 测试 PPTX
	t.Run("PPTX文件读取", func(t *testing.T) {
		pptxPath := filepath.Join(testDir, "test.pptx")
		if _, err := os.Stat(pptxPath); err == nil {
			doc, err := ReadDocument(pptxPath)
			if err != nil {
				t.Errorf("读取 PPTX 失败: %v", err)
				return
			}
			if doc.Content == "" {
				t.Error("PPTX 内容为空")
			}

			// 测试获取幻灯片
			reader := &PptxReader{}
			slides, err := reader.GetSlides(pptxPath)
			if err != nil {
				t.Errorf("获取幻灯片失败: %v", err)
			} else {
				t.Logf("PPTX 幻灯片数量: %d", len(slides))
			}
		} else {
			t.Skip("测试文件不存在: test.pptx")
		}
	})

	// 测试 TXT
	t.Run("TXT文件读取", func(t *testing.T) {
		txtPath := filepath.Join(testDir, "test.txt")
		if _, err := os.Stat(txtPath); err == nil {
			doc, err := ReadDocument(txtPath)
			if err != nil {
				t.Errorf("读取 TXT 失败: %v", err)
				return
			}
			t.Logf("TXT 内容长度: %d 字符", len(doc.Content))
		} else {
			t.Skip("测试文件不存在: test.txt")
		}
	})

	// 测试 CSV
	t.Run("CSV文件读取", func(t *testing.T) {
		csvPath := filepath.Join(testDir, "test.csv")
		if _, err := os.Stat(csvPath); err == nil {
			doc, err := ReadDocument(csvPath)
			if err != nil {
				t.Errorf("读取 CSV 失败: %v", err)
				return
			}
			t.Logf("CSV 内容长度: %d 字符", len(doc.Content))

			// 测试获取记录
			reader := &CsvReader{}
			records, err := reader.GetRecords(csvPath)
			if err != nil {
				t.Errorf("获取 CSV 记录失败: %v", err)
			} else {
				t.Logf("CSV 记录数量: %d 行", len(records))
			}
		} else {
			t.Skip("测试文件不存在: test.csv")
		}
	})

	// 测试 Markdown
	t.Run("Markdown文件读取", func(t *testing.T) {
		mdPath := filepath.Join(testDir, "test.md")
		if _, err := os.Stat(mdPath); err == nil {
			doc, err := ReadDocument(mdPath)
			if err != nil {
				t.Errorf("读取 Markdown 失败: %v", err)
				return
			}
			t.Logf("Markdown 内容长度: %d 字符", len(doc.Content))
		} else {
			t.Skip("测试文件不存在: test.md")
		}
	})
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

// BenchmarkReadDocument 性能基准测试
func BenchmarkReadDocument(b *testing.B) {
	testFiles := []string{
		"testdata/test.docx",
		"testdata/test.txt",
		"testdata/test.csv",
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
