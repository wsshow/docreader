package docreader

import (
	"errors"
	"fmt"
)

// 预定义的错误类型
var (
	// ErrUnsupportedFormat 不支持的文件格式
	ErrUnsupportedFormat = errors.New("unsupported file format")

	// ErrFileNotFound 文件不存在
	ErrFileNotFound = errors.New("file not found")

	// ErrFileOpen 无法打开文件
	ErrFileOpen = errors.New("failed to open file")

	// ErrFileRead 读取文件失败
	ErrFileRead = errors.New("failed to read file")

	// ErrFileParse 解析文件失败
	ErrFileParse = errors.New("failed to parse file")

	// ErrInvalidFormat 文件格式无效
	ErrInvalidFormat = errors.New("invalid file format")

	// ErrEmptyFile 文件为空
	ErrEmptyFile = errors.New("file is empty")

	// ErrSheetNotFound 工作表不存在
	ErrSheetNotFound = errors.New("sheet not found")
)

// DocumentError 文档错误结构
type DocumentError struct {
	Op       string // 操作名称
	FilePath string // 文件路径
	Err      error  // 原始错误
}

// Error 实现 error 接口
func (e *DocumentError) Error() string {
	if e.FilePath != "" {
		return fmt.Sprintf("%s: %s: %v", e.Op, e.FilePath, e.Err)
	}
	return fmt.Sprintf("%s: %v", e.Op, e.Err)
}

// Unwrap 返回原始错误，支持 errors.Is 和 errors.As
func (e *DocumentError) Unwrap() error {
	return e.Err
}

// NewError 创建新的文档错误
func NewError(op, filePath string, err error) error {
	return &DocumentError{
		Op:       op,
		FilePath: filePath,
		Err:      err,
	}
}

// WrapError 包装错误并添加上下文信息
func WrapError(op, filePath string, err error) error {
	if err == nil {
		return nil
	}
	return NewError(op, filePath, err)
}

// IsUnsupportedFormat 检查是否为不支持的格式错误
func IsUnsupportedFormat(err error) bool {
	return errors.Is(err, ErrUnsupportedFormat)
}

// IsFileNotFound 检查是否为文件不存在错误
func IsFileNotFound(err error) bool {
	return errors.Is(err, ErrFileNotFound)
}

// IsFileOpen 检查是否为文件打开错误
func IsFileOpen(err error) bool {
	return errors.Is(err, ErrFileOpen)
}

// IsFileRead 检查是否为文件读取错误
func IsFileRead(err error) bool {
	return errors.Is(err, ErrFileRead)
}

// IsFileParse 检查是否为文件解析错误
func IsFileParse(err error) bool {
	return errors.Is(err, ErrFileParse)
}
