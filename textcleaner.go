package docreader

import (
	"regexp"
	"strings"
	"unicode"
)

// TextCleaner 提供文本清理功能，用于优化大模型理解
type TextCleaner struct {
	// 是否移除多余空行（连续的空行压缩为一个）
	RemoveExtraBlankLines bool
	// 是否移除行首行尾空格
	TrimSpaces bool
	// 是否移除多余空格（连续空格压缩为一个）
	RemoveExtraSpaces bool
	// 是否移除特殊控制字符
	RemoveControlChars bool
	// 最大连续空行数（0表示不限制）
	MaxConsecutiveBlankLines int
}

// DefaultTextCleaner 返回默认配置的文本清理器
func DefaultTextCleaner() *TextCleaner {
	return &TextCleaner{
		RemoveExtraBlankLines:    true,
		TrimSpaces:               true,
		RemoveExtraSpaces:        true,
		RemoveControlChars:       true,
		MaxConsecutiveBlankLines: 1,
	}
}

// Clean 清理文本内容
func (tc *TextCleaner) Clean(text string) string {
	if text == "" {
		return ""
	}

	// 1. 移除控制字符（保留换行符、制表符等有意义的空白字符）
	if tc.RemoveControlChars {
		text = tc.removeControlChars(text)
	}

	// 2. 统一换行符为 \n
	text = normalizeLineBreaks(text)

	// 3. 按行处理
	lines := strings.Split(text, "\n")
	var cleanedLines []string
	consecutiveBlankLines := 0

	for _, line := range lines {
		// 移除行首行尾空格
		if tc.TrimSpaces {
			line = strings.TrimSpace(line)
		}

		// 移除多余空格
		if tc.RemoveExtraSpaces && line != "" {
			line = tc.removeExtraSpaces(line)
		}

		// 处理空行
		if line == "" {
			consecutiveBlankLines++

			// 如果 MaxConsecutiveBlankLines 为 0，移除所有空行
			if tc.MaxConsecutiveBlankLines == 0 {
				continue
			}

			// 如果设置了最大连续空行数（大于0）
			if tc.MaxConsecutiveBlankLines > 0 && consecutiveBlankLines > tc.MaxConsecutiveBlankLines {
				continue // 跳过多余的空行
			}

			// 如果要移除多余空行（且 MaxConsecutiveBlankLines < 0 表示不限制）
			if tc.RemoveExtraBlankLines && consecutiveBlankLines > 1 && tc.MaxConsecutiveBlankLines != -1 {
				continue // 跳过多余的空行
			}
		} else {
			consecutiveBlankLines = 0
		}

		cleanedLines = append(cleanedLines, line)
	}

	// 4. 移除开头和结尾的空行
	cleanedLines = trimEmptyLines(cleanedLines)

	// 5. 合并结果
	result := strings.Join(cleanedLines, "\n")

	return result
}

// removeControlChars 移除控制字符，保留必要的空白字符
func (tc *TextCleaner) removeControlChars(text string) string {
	var builder strings.Builder
	builder.Grow(len(text))

	for _, r := range text {
		// 保留可打印字符、换行符、制表符、回车符
		if unicode.IsPrint(r) || r == '\n' || r == '\t' || r == '\r' {
			builder.WriteRune(r)
		}
	}

	return builder.String()
}

// removeExtraSpaces 移除多余的空格（连续空格压缩为一个）
func (tc *TextCleaner) removeExtraSpaces(text string) string {
	// 使用正则表达式将多个连续空格替换为单个空格
	re := regexp.MustCompile(`[ \t]+`)
	return re.ReplaceAllString(text, " ")
}

// normalizeLineBreaks 统一换行符
func normalizeLineBreaks(text string) string {
	// 将 \r\n 和 \r 统一替换为 \n
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	return text
}

// trimEmptyLines 移除开头和结尾的空行
func trimEmptyLines(lines []string) []string {
	// 移除开头的空行
	start := 0
	for start < len(lines) && lines[start] == "" {
		start++
	}

	// 移除结尾的空行
	end := len(lines)
	for end > start && lines[end-1] == "" {
		end--
	}

	if start >= end {
		return []string{}
	}

	return lines[start:end]
}

// CleanText 使用默认配置清理文本的便捷函数
func CleanText(text string) string {
	cleaner := DefaultTextCleaner()
	return cleaner.Clean(text)
}

// CleanTextMinimal 使用最小清理配置（仅清理基本的空白）
func CleanTextMinimal(text string) string {
	cleaner := &TextCleaner{
		RemoveExtraBlankLines:    false,
		TrimSpaces:               true,
		RemoveExtraSpaces:        true,
		RemoveControlChars:       true,
		MaxConsecutiveBlankLines: -1, // 不限制空行数
	}
	return cleaner.Clean(text)
}

// CleanTextAggressive 使用激进的清理配置（最大程度压缩空间）
func CleanTextAggressive(text string) string {
	cleaner := &TextCleaner{
		RemoveExtraBlankLines:    true,
		TrimSpaces:               true,
		RemoveExtraSpaces:        true,
		RemoveControlChars:       true,
		MaxConsecutiveBlankLines: 0, // 不允许空行
	}
	return cleaner.Clean(text)
}
