# DocReader - Go 文档读取库

一个用于读取常见文档格式的 Go 语言库，支持 Office 文档、PDF、文本文件等多种格式。

## 功能特性

### Office 文档

- ✅ 读取 **DOCX** (Word 文档) 的文本内容和元数据
- ✅ 读取 **XLSX** (Excel 表格) 的文本内容和结构化数据
- ✅ 读取 **PPTX** (PowerPoint 演示文稿) 的文本内容

### PDF 文档

- ✅ 读取 **PDF** 文件的文本内容和元数据

### 文本格式

- ✅ 读取 **TXT** 纯文本文件
- ✅ 读取 **CSV** 表格文件（支持结构化数据）
- ✅ 读取 **Markdown** (.md) 文件
- ✅ 读取 **RTF** 富文本格式（基础文本提取）

### 其他特性

- ✅ 统一的接口设计，自动识别文件格式
- ✅ 提取文档元数据（标题、作者、创建时间等）
- ✅ 支持中文内容

## 安装

```bash
go get github.com/wsshow/docreader
```

## 快速开始

### 基本用法 - 自动识别文件类型

```go
package main

import (
    "fmt"
    "log"
    "github.com/wsshow/docreader"
)

func main() {
    // 自动识别文件格式并读取
    doc, err := docreader.ReadDocument("example.docx")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("文件: %s\n", doc.FilePath)
    fmt.Printf("内容: %s\n", doc.Content)
    fmt.Printf("元数据: %v\n", doc.Metadata)
}
```

### DOCX - Word 文档

```go
// 方式 1: 使用统一接口
doc, err := docreader.ReadDocument("document.docx")
if err != nil {
    log.Fatal(err)
}
fmt.Println(doc.Content)

// 方式 2: 使用专用读取器
reader := &docreader.DocxReader{}
content, err := reader.ReadText("document.docx")
if err != nil {
    log.Fatal(err)
}
fmt.Println(content)

// 获取元数据
metadata, err := reader.GetMetadata("document.docx")
fmt.Printf("标题: %s\n", metadata["title"])
fmt.Printf("作者: %s\n", metadata["creator"])
```

### PDF 文件

```go
// 读取 PDF 内容
doc, err := docreader.ReadDocument("document.pdf")
if err != nil {
    log.Fatal(err)
}
fmt.Println(doc.Content)

// 获取 PDF 元数据
reader := &docreader.PdfReader{}
metadata, err := reader.GetMetadata("document.pdf")
fmt.Printf("页数: %s\n", metadata["pages"])
fmt.Printf("作者: %s\n", metadata["author"])
```

### XLSX - Excel 表格

```go
// 基本读取
doc, err := docreader.ReadDocument("spreadsheet.xlsx")
if err != nil {
    log.Fatal(err)
}
fmt.Println(doc.Content)

// 高级用法 - 获取结构化数据
reader := &docreader.XlsxReader{}

// 获取指定工作表的数据
rows, err := reader.GetSheetData("spreadsheet.xlsx", "Sheet1")
if err != nil {
    log.Fatal(err)
}
for _, row := range rows {
    fmt.Println(row) // []string
}

// 获取所有工作表的数据
allSheets, err := reader.GetAllSheetsData("spreadsheet.xlsx")
for sheetName, rows := range allSheets {
    fmt.Printf("工作表: %s\n", sheetName)
    for _, row := range rows {
        fmt.Println(row)
    }
}
```

### PPTX - PowerPoint 演示文稿

```go
// 基本读取
doc, err := docreader.ReadDocument("presentation.pptx")
if err != nil {
    log.Fatal(err)
}
fmt.Println(doc.Content)

// 高级用法 - 按幻灯片分组
reader := &docreader.PptxReader{}
slides, err := reader.GetSlides("presentation.pptx")
if err != nil {
    log.Fatal(err)
}

for i, slide := range slides {
    fmt.Printf("=== 幻灯片 %d ===\n", i+1)
    fmt.Println(slide)
}

// 获取元数据
metadata, err := reader.GetMetadata("presentation.pptx")
fmt.Printf("幻灯片总数: %s\n", metadata["slide_count"])
```

### TXT - 纯文本文件

```go
// 读取文本文件
doc, err := docreader.ReadDocument("document.txt")
if err != nil {
    log.Fatal(err)
}
fmt.Println(doc.Content)
```

### CSV - 表格文件

```go
// 基本读取
doc, err := docreader.ReadDocument("data.csv")
if err != nil {
    log.Fatal(err)
}
fmt.Println(doc.Content)

// 高级用法 - 获取结构化数据
reader := &docreader.CsvReader{}
records, err := reader.GetRecords("data.csv")
if err != nil {
    log.Fatal(err)
}

for _, row := range records {
    fmt.Println(row) // []string
}
```

### Markdown - Markdown 文件

```go
// 读取 Markdown 文件
doc, err := docreader.ReadDocument("README.md")
if err != nil {
    log.Fatal(err)
}
fmt.Println(doc.Content)
```

### RTF - 富文本格式

```go
// 读取 RTF 文件（基础文本提取）
doc, err := docreader.ReadDocument("document.rtf")
if err != nil {
    log.Fatal(err)
}
fmt.Println(doc.Content)
```

## 高级配置

### 精确控制读取内容

DocReader 提供了强大的配置系统，允许你精确控制要读取的页面和行。

#### 基本配置示例

```go
// 示例 1: 读取指定页码
config := docreader.NewReadConfig().
    WithPages(0, 2, 4) // 读取第 0、2、4 页

result, err := docreader.ReadDocumentWithConfig("document.pdf", config)

// 示例 2: 读取页码范围
config := docreader.NewReadConfig().
    WithPageRange(0, 5) // 读取第 0-5 页

// 示例 3: 读取指定行
config := docreader.NewReadConfig().
    WithLines(0, 5, 10) // 每页只读取第 0、5、10 行

// 示例 4: 读取行范围
config := docreader.NewReadConfig().
    WithLineRange(5, 15) // 每页只读取第 5-15 行

// 示例 5: 组合多个范围
config := docreader.NewReadConfig().
    WithPageRange(0, 2).
    WithPageRange(5, 7).  // 读取第 0-2 页和第 5-7 页
    WithLineRange(0, 10)  // 每页读取第 0-10 行
```

#### 为不同页面设置不同的行配置

```go
// 为特定页面配置不同的行
config := docreader.NewReadConfig().
    AddPageLines(0, 1, 3, 5).           // 第 0 页：读取第 1、3、5 行
    AddPageLineRange(1, 0, 10).         // 第 1 页：读取第 0-10 行
    AddPageConfig(2, []int{2, 4}, [][2]int{{6, 8}}) // 第 2 页：读取第 2、4 行和第 6-8 行

result, err := docreader.ReadDocumentWithConfig("document.pdf", config)
```

#### XLSX 工作表筛选

```go
// 只读取指定的工作表
config := docreader.NewReadConfig().
    WithSheetNames("Sheet1", "Summary")

result, err := docreader.ReadDocumentWithConfig("spreadsheet.xlsx", config)
```

#### 处理结构化结果

```go
result, err := docreader.ReadDocumentWithConfig("document.pdf", config)
if err != nil {
    log.Fatal(err)
}

// 访问结构化数据
fmt.Printf("总页数: %d\n", result.TotalPages)
fmt.Printf("总行数: %d\n", result.TotalLines)

// 遍历每一页
for _, page := range result.Pages {
    fmt.Printf("页码: %d, 行数: %d\n", page.PageNumber, page.TotalLines)
    for i, line := range page.Lines {
        fmt.Printf("  行 %d: %s\n", i, line)
    }
}

// 获取完整内容
fmt.Println(result.Content)
```

### 文本清理

DocReader 提供了智能的文本清理功能，可以优化提取的文本内容，特别适合用于大模型处理。

#### 使用预设清理配置

```go
// 使用默认清理（推荐）
doc, err := docreader.ReadDocumentWithClean("document.pdf")
// 效果：移除行首尾空格、压缩多余空格、移除控制字符、最多保留1个连续空行

// 最小清理（保留格式）
doc, err := docreader.ReadDocument("document.txt")
doc.CleanContentMinimal()
// 效果：仅移除行首尾空格、压缩多余空格、移除控制字符，保留所有空行

// 激进清理（最大压缩）
doc, err := docreader.ReadDocument("document.docx")
doc.CleanContentAggressive()
// 效果：移除所有空行，压缩空格，移除控制字符
```

#### 自定义清理配置

```go
// 创建自定义清理器
cleaner := &docreader.TextCleaner{
    TrimSpaces:         true,  // 移除行首行尾空格
    RemoveExtraSpaces:  true,  // 压缩连续空格为一个
    RemoveControlChars: true,  // 移除控制字符
    MaxBlankLines:      2,     // 最多保留 2 个连续空行（-1=不限制，0=移除所有）
}

// 在读取时应用清理
doc, err := docreader.ReadDocumentWithCleanConfig("document.pdf", cleaner)

// 或者对已读取的文档应用清理
doc, err := docreader.ReadDocument("document.txt")
doc.CleanContentWith(cleaner)
```

#### TextCleaner 配置说明

```go
type TextCleaner struct {
    // TrimSpaces: 是否移除行首行尾空格
    TrimSpaces bool

    // RemoveExtraSpaces: 是否将连续空格压缩为一个
    RemoveExtraSpaces bool

    // RemoveControlChars: 是否移除特殊控制字符（保留换行符和制表符）
    RemoveControlChars bool

    // MaxBlankLines: 最大连续空行数
    //   -1: 不限制，保留所有空行
    //    0: 移除所有空行
    //    1: 最多保留 1 个连续空行（压缩多余空行）
    //    N: 最多保留 N 个连续空行
    MaxBlankLines int
}
```

## API 文档

### 核心接口

#### `ReadDocument(filePath string) (*Document, error)`

自动识别文件格式并读取内容，返回包含内容和元数据的 Document 对象。

#### `ReadDocumentWithClean(filePath string) (*Document, error)`

读取文档并自动应用默认文本清理。

#### `ReadDocumentWithCleanConfig(filePath string, cleaner *TextCleaner) (*Document, error)`

读取文档并应用自定义文本清理配置。

#### `ReadDocumentWithConfig(filePath string, config *ReadConfig) (*DocumentResult, error)`

根据配置精确读取文档，返回结构化的结果。

#### `NewReadConfig() *ReadConfig`

创建一个新的读取配置对象，支持链式调用。

#### `DocumentReader` 接口

所有读取器都实现此接口：

- `ReadText(filePath string) (string, error)` - 读取文本内容
- `GetMetadata(filePath string) (map[string]string, error)` - 获取元数据

#### `ConfigurableReader` 接口

支持高级配置的读取器接口（所有读取器都实现）：

- `ReadWithConfig(filePath string, config *ReadConfig) (*DocumentResult, error)` - 根据配置读取文档

### 配置结构

#### ReadConfig 配置方法

```go
// 页面选择
config.WithPages(pages ...int)              // 设置要读取的离散页码
config.WithPageRange(start, end int)        // 添加页码范围

// 全局行选择（应用到所有页）
config.WithLines(lines ...int)              // 设置要读取的离散行号
config.WithLineRange(start, end int)        // 添加行号范围

// 页面级行配置（覆盖全局配置）
config.AddPageLines(pageIndex int, lines ...int)                          // 为指定页添加离散行
config.AddPageLineRange(pageIndex, start, end int)                       // 为指定页添加行范围
config.AddPageConfig(pageIndex int, lineIndexes []int, lineRanges [][2]int) // 为指定页添加完整配置

// XLSX 特有
config.WithSheetNames(names ...string)      // 设置要读取的工作表名称
```

#### 核心数据结构

```go
// Selector 统一的选择器，用于选择页码或行号
type Selector struct {
    Indexes []int       // 离散索引：[0, 2, 5]
    Ranges  [][2]int    // 连续范围：[[0,2], [5,10]]
}

// ReadConfig 读取配置
type ReadConfig struct {
    PageSelector Selector      // 页面选择器
    LineSelector Selector      // 全局行选择器
    PageConfigs  []PageConfig  // 页面级配置（优先级高于全局）
    SheetNames   []string      // XLSX 工作表名称
}

// DocumentResult 结构化的文档读取结果
type DocumentResult struct {
    FilePath   string
    Pages      []PageContent      // 结构化页面内容
    TotalPages int
    TotalLines int
    Metadata   map[string]string
    Content    string             // 完整文本内容
}

// PageContent 单页内容
type PageContent struct {
    PageNumber int
    PageName   string   // 工作表名称（XLSX）
    Lines      []string
    TotalLines int
}
```

### 专用读取器

#### DocxReader

- `ReadText()` - 读取段落和表格文本
- `GetMetadata()` - 获取标题、作者、创建/修改时间等

#### PdfReader

- `ReadText()` - 逐页读取文本内容
- `GetMetadata()` - 获取页数、作者、创建时间等

#### XlsxReader

- `ReadText()` - 读取所有工作表的文本
- `GetMetadata()` - 获取工作表列表、文档属性等
- `GetSheetData(filePath, sheetName string)` - 获取指定工作表的结构化数据
- `GetAllSheetsData(filePath string)` - 获取所有工作表的结构化数据

#### PptxReader

- `ReadText()` - 读取所有幻灯片的文本
- `GetMetadata()` - 获取幻灯片数量、标题等
- `GetSlides(filePath string)` - 按幻灯片分组获取文本

#### TxtReader

- `ReadText()` - 读取纯文本内容
- `GetMetadata()` - 获取文件大小、修改时间等

#### CsvReader

- `ReadText()` - 读取 CSV 文件的格式化文本
- `GetMetadata()` - 获取行数、列数、文件信息等
- `GetRecords(filePath string)` - 获取结构化的 CSV 数据

#### MdReader

- `ReadText()` - 读取 Markdown 原始内容
- `GetMetadata()` - 获取文件大小、修改时间等

#### RtfReader

- `ReadText()` - 提取 RTF 文件的纯文本内容
- `GetMetadata()` - 获取文件大小、修改时间等

## 支持的元数据

### DOCX/PPTX

- title - 标题
- subject - 主题
- creator - 创建者
- description - 描述
- created - 创建时间
- modified - 修改时间

### PDF

- title - 标题
- author - 作者
- subject - 主题
- creator - 创建程序
- producer - 生成程序
- creation_date - 创建日期
- modification_date - 修改日期
- pages - 页数

### XLSX

- title - 标题
- subject - 主题
- creator - 创建者
- sheets - 工作表列表
- sheet_count - 工作表数量
- active_sheet - 活动工作表

## 已知限制

### PDF 中文字符支持

当前 PDF 读取器使用 `ledongthuc/pdf` 库，该库对某些 PDF 文件中的中文字符（CJK 字体）支持有限。如果 PDF 文件使用了嵌入式中文字体或特殊编码，可能会出现乱码。

**建议**：

- 对于包含大量中文内容的 PDF，建议使用其他专业 PDF 处理工具
- 英文和数字内容可以正常提取
- 元数据提取不受影响

## 错误处理

库提供了统一的错误封装和类型检查功能，方便进行精确的错误处理。

### 预定义错误类型

```go
var (
    ErrUnsupportedFormat = errors.New("unsupported file format")  // 不支持的文件格式
    ErrFileNotFound      = errors.New("file not found")           // 文件不存在
    ErrFileOpen          = errors.New("failed to open file")      // 无法打开文件
    ErrFileRead          = errors.New("failed to read file")      // 读取文件失败
    ErrFileParse         = errors.New("failed to parse file")     // 解析文件失败
    ErrInvalidFormat     = errors.New("invalid file format")      // 文件格式无效
    ErrEmptyFile         = errors.New("file is empty")            // 文件为空
    ErrSheetNotFound     = errors.New("sheet not found")          // 工作表不存在
)
```

### 基本错误处理

```go
doc, err := docreader.ReadDocument("file.docx")
if err != nil {
    log.Printf("读取文件失败: %v", err)
    return
}
```

### 使用错误类型检查

```go
doc, err := docreader.ReadDocument("file.unknown")
if err != nil {
    // 使用 errors.Is 进行错误类型判断
    if errors.Is(err, docreader.ErrUnsupportedFormat) {
        log.Println("不支持的文件格式")
    } else if errors.Is(err, docreader.ErrFileNotFound) {
        log.Println("文件不存在")
    } else if errors.Is(err, docreader.ErrFileOpen) {
        log.Println("无法打开文件")
    } else {
        log.Printf("其他错误: %v", err)
    }
    return
}
```

### 使用辅助函数

```go
doc, err := docreader.ReadDocument("file.pdf")
if err != nil {
    // 使用辅助函数进行错误检查
    if docreader.IsUnsupportedFormat(err) {
        log.Println("不支持的文件格式")
    } else if docreader.IsFileNotFound(err) {
        log.Println("文件不存在")
    } else if docreader.IsFileOpen(err) {
        log.Println("无法打开文件")
    } else if docreader.IsFileRead(err) {
        log.Println("读取文件失败")
    } else if docreader.IsFileParse(err) {
        log.Println("解析文件失败")
    } else {
        log.Printf("未知错误: %v", err)
    }
    return
}
```

### 获取详细错误信息

```go
doc, err := docreader.ReadDocument("file.docx")
if err != nil {
    // 错误信息包含操作名称和文件路径
    // 格式: "操作名称: 文件路径: 错误详情"
    log.Printf("详细错误: %v", err)

    // 使用 errors.Unwrap 获取原始错误
    if unwrapped := errors.Unwrap(err); unwrapped != nil {
        log.Printf("原始错误: %v", unwrapped)
    }
    return
}
```

### 完整示例

```go
package main

import (
    "errors"
    "log"
    "github.com/wsshow/docreader"
)

func main() {
    filePath := "document.pdf"

    doc, err := docreader.ReadDocument(filePath)
    if err != nil {
        handleError(err)
        return
    }

    log.Printf("成功读取文档，内容长度: %d", len(doc.Content))
}

func handleError(err error) {
    switch {
    case docreader.IsUnsupportedFormat(err):
        log.Println("错误: 不支持的文件格式，请使用 .docx, .pdf, .xlsx, .pptx, .txt, .csv, .md 或 .rtf 格式")
    case docreader.IsFileNotFound(err):
        log.Println("错误: 文件不存在，请检查文件路径")
    case docreader.IsFileOpen(err):
        log.Println("错误: 无法打开文件，请检查文件权限")
    case docreader.IsFileRead(err):
        log.Println("错误: 读取文件失败，文件可能已损坏")
    case docreader.IsFileParse(err):
        log.Println("错误: 解析文件失败，文件格式可能不正确")
    default:
        log.Printf("错误: %v", err)
    }
}
```
