package flashtext

import (
	"fmt"
	"testing"
)

// // 基础功能测试
// Match:  his
// 1
// 4
// =========
// Match:  is
// 2
// 4
// =========
// Match:  she
// 3
// 6
// =========
// Match:  he
// 4
// 6
// =========
// Match:  hers
// 4
// 8
// =========
// Match:  she
// 7
// 10
// =========
// Match:  he
// 8
// 10
// =========
// Match:  share
// 10
// 15
// =========
func TestBasicMatching(t *testing.T) {
	kp := NewKeywordProcessor() // case insensitive

	// 2. 添加关键词
	kp.AddKeyWord("Big Data")
	kp.AddKeyWord("Python")
	kp.Build()
	// 3. 匹配
	text := "I love Big Data and Python."
	matches := kp.ExtractKeywords(text)
	fmt.Println(matches)
}

// 中文测试
func TestChineseMatching(t *testing.T) {
	kp := NewKeywordProcessor()
	kp.AddKeyWord("毛泽东").Build()

	text := "ahisHershare毛泽东dsadsa"
	matches := kp.ExtractKeywords(text)
	fmt.Println(matches)
	if len(matches) != 1 {
		t.Fatalf("期望匹配1个关键词, 实际 %d 个", len(matches))
	}

	match := matches[0]
	if match.MatchString() != "毛泽东" {
		t.Errorf("期望匹配 '毛泽东', 实际 '%s'", match.MatchString())
	}

	// 验证字节位置
	if text[match.Start():match.End()] != "毛泽东" {
		t.Errorf("字节位置不正确")
	}
}

// 大小写不敏感测试
func TestCaseInsensitive(t *testing.T) {
	kp := NewKeywordProcessor() // 不区分大小写
	kp.AddKeyWord("apple").Build()

	tests := []struct {
		text     string
		expected int
	}{
		{"I have an apple", 1},
		{"I have an Apple", 1},
		{"I have an APPLE", 1},
		{"I have an aPpLe", 1},
		{"no fruit here", 0},
	}

	for _, tt := range tests {
		matches := kp.ExtractKeywords(tt.text)
		if len(matches) != tt.expected {
			t.Errorf("文本 '%s': 期望 %d 个匹配, 实际 %d 个",
				tt.text, tt.expected, len(matches))
		}
	}
}

// 大小写敏感测试
func TestCaseSensitive(t *testing.T) {
	kp := NewKeywordProcessor(WithCaseSensitive()) // 区分大小写
	kp.AddKeyWord("Apple").Build()

	tests := []struct {
		text     string
		expected int
	}{
		{"I have an Apple", 1},
		{"I have an apple", 0}, // 小写不匹配
		{"I have an APPLE", 0}, // 大写不匹配
	}

	for _, tt := range tests {
		matches := kp.ExtractKeywords(tt.text)
		if len(matches) != tt.expected {
			t.Errorf("文本 '%s': 期望 %d 个匹配, 实际 %d 个",
				tt.text, tt.expected, len(matches))
		}
	}
}

// 边缘情况测试
func TestEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		keywords []string
		text     string
		expected int
	}{
		{"空文本", []string{"test"}, "", 0},
		{"空关键词", []string{}, "test text", 0},
		{"空关键词字符串", []string{""}, "test", 0},
		{"完全匹配", []string{"hello"}, "hello", 1},
		{"不匹配", []string{"hello"}, "world", 0},
		{"重叠匹配", []string{"he", "she", "hers"}, "hershey", 4}, // he, she, hers, he
		{"单字符", []string{"a"}, "a b c a", 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kp := NewKeywordProcessor()
			kp.AddKeywordsFromList(tt.keywords).Build()
			matches := kp.ExtractKeywords(tt.text)
			if len(matches) != tt.expected {
				t.Errorf("期望 %d 个匹配, 实际 %d 个", tt.expected, len(matches))
			}
		})
	}
}

// Bytes方法测试
func TestExtractFromBytes(t *testing.T) {
	kp := NewKeywordProcessor()
	kp.AddKeywordsFromList([]string{"hello", "world"}).Build()

	text := []byte("hello world hello")
	matches := kp.ExtractKeywordsFromBytes(text)

	if len(matches) != 3 {
		t.Errorf("期望3个匹配, 实际 %d 个", len(matches))
	}
}
