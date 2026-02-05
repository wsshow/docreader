package docreader

import (
	"strings"
	"testing"
)

func TestTextCleaner(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "移除多余空格",
			input:    "Hello    world   test",
			expected: "Hello world test",
		},
		{
			name:     "移除多余空行",
			input:    "Line 1\n\n\n\nLine 2\n\n\n\nLine 3",
			expected: "Line 1\n\nLine 2\n\nLine 3",
		},
		{
			name:     "移除行首行尾空格",
			input:    "  Line 1  \n  Line 2  \n  Line 3  ",
			expected: "Line 1\nLine 2\nLine 3",
		},
		{
			name:     "统一换行符",
			input:    "Line 1\r\nLine 2\rLine 3\nLine 4",
			expected: "Line 1\nLine 2\nLine 3\nLine 4",
		},
		{
			name:     "移除开头和结尾的空行",
			input:    "\n\n\nLine 1\nLine 2\n\n\n",
			expected: "Line 1\nLine 2",
		},
		{
			name:     "综合测试",
			input:    "\n\n  Line 1   with   spaces  \n\n\n\n  Line 2  \n\n  Line 3  \n\n\n",
			expected: "Line 1 with spaces\n\nLine 2\n\nLine 3",
		},
		{
			name:     "空字符串",
			input:    "",
			expected: "",
		},
		{
			name:     "只有空格和空行",
			input:    "   \n\n   \n   ",
			expected: "",
		},
	}

	cleaner := DefaultTextCleaner()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleaner.Clean(tt.input)
			if result != tt.expected {
				t.Errorf("期望:\n%q\n实际:\n%q", tt.expected, result)
			}
		})
	}
}

func TestCleanTextMinimal(t *testing.T) {
	input := "Line 1\n\n\n\nLine 2\n\n\n\nLine 3"
	result := CleanTextMinimal(input)

	// 最小清理不应移除多余空行
	if !strings.Contains(result, "\n\n\n") {
		t.Error("最小清理模式应保留多余空行")
	}
}

func TestCleanTextAggressive(t *testing.T) {
	input := "Line 1\n\n\n\nLine 2\n\n\n\nLine 3"
	result := CleanTextAggressive(input)

	// 激进清理应该移除所有空行
	if strings.Contains(result, "\n\n") {
		t.Error("激进清理模式应移除所有空行")
	}

	expected := "Line 1\nLine 2\nLine 3"
	if result != expected {
		t.Errorf("激进清理结果不符合预期\n期望: %q\n实际: %q", expected, result)
	}
}

func TestTextCleanerCustomConfig(t *testing.T) {
	cleaner := &TextCleaner{
		TrimSpaces:         true,
		RemoveExtraSpaces:  false, // 不移除多余空格
		RemoveControlChars: true,
		MaxBlankLines:      2, // 允许最多2个连续空行
	}

	input := "Hello    world\n\n\n\n\nTest"
	result := cleaner.Clean(input)

	// 应该保留多余空格
	if !strings.Contains(result, "    ") {
		t.Error("应该保留多余空格")
	}

	// 应该限制连续空行数量
	lines := strings.Split(result, "\n")
	consecutiveEmpty := 0
	maxConsecutive := 0

	for _, line := range lines {
		if line == "" {
			consecutiveEmpty++
			if consecutiveEmpty > maxConsecutive {
				maxConsecutive = consecutiveEmpty
			}
		} else {
			consecutiveEmpty = 0
		}
	}

	if maxConsecutive > 2 {
		t.Errorf("连续空行数量应不超过2，实际为: %d", maxConsecutive)
	}
}

func TestRemoveControlChars(t *testing.T) {
	cleaner := DefaultTextCleaner()

	// 包含控制字符的文本
	input := "Hello\x00\x01\x02World\x03\x04Test"
	result := cleaner.Clean(input)

	// 应该移除控制字符
	for _, r := range result {
		if r < 32 && r != '\n' && r != '\t' {
			t.Errorf("不应包含控制字符: %v", r)
		}
	}
}

func TestNormalizeLineBreaks(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Windows换行符",
			input:    "Line1\r\nLine2\r\nLine3",
			expected: "Line1\nLine2\nLine3",
		},
		{
			name:     "Mac换行符",
			input:    "Line1\rLine2\rLine3",
			expected: "Line1\nLine2\nLine3",
		},
		{
			name:     "Unix换行符",
			input:    "Line1\nLine2\nLine3",
			expected: "Line1\nLine2\nLine3",
		},
		{
			name:     "混合换行符",
			input:    "Line1\r\nLine2\rLine3\nLine4",
			expected: "Line1\nLine2\nLine3\nLine4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeLineBreaks(tt.input)
			if result != tt.expected {
				t.Errorf("期望: %q, 实际: %q", tt.expected, result)
			}
		})
	}
}

func BenchmarkCleanText(b *testing.B) {
	input := strings.Repeat("Line 1   with   spaces  \n\n\n", 100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CleanText(input)
	}
}

func BenchmarkCleanTextAggressive(b *testing.B) {
	input := strings.Repeat("Line 1   with   spaces  \n\n\n", 100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CleanTextAggressive(input)
	}
}
