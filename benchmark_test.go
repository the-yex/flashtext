package flashtext

import "testing"

// 性能基准测试

// 基础性能测试 - 中文文本
func BenchmarkExtractKeywords(b *testing.B) {
	keywords := []string{
		"A", "AV", "AV演员", "无名氏", "AV演员色情", "日本AV女优",
	}
	text := "日本AV演员兼电视、电影演员。无名氏AV女优是xx出道, 日本AV女优们最精彩的表演是AV演员色情表演"

	kp := NewKeywordProcessor()
	kp.AddKeywordsFromList(keywords).Build()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		kp.ExtractKeywords(text)
	}
}

func BenchmarkExtractFromBytes(b *testing.B) {
	keywords := []string{
		"A", "AV", "AV演员", "无名氏", "AV演员色情", "日本AV女优",
	}
	text := []byte("日本AV演员兼电视、电影演员。无名氏AV女优是xx出道, 日本AV女优们最精彩的表演是AV演员色情表演")

	kp := NewKeywordProcessor()
	kp.AddKeywordsFromList(keywords).Build()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		kp.ExtractKeywordsFromBytes(text)
	}
}

// 测试不同关键词数量的性能影响

// 小关键词集 (10个)
func BenchmarkSmallKeywordSet(b *testing.B) {
	keywords := []string{
		"hello", "world", "test", "golang", "code",
		"data", "file", "user", "name", "text",
	}
	text := "This is a test text with hello world and some golang code for user data processing."

	kp := NewKeywordProcessor()
	kp.AddKeywordsFromList(keywords).Build()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		kp.ExtractKeywords(text)
	}
}

// 中等关键词集 (100个)
func BenchmarkMediumKeywordSet(b *testing.B) {
	// 生成100个关键词
	keywords := make([]string, 100)
	for i := 0; i < 100; i++ {
		keywords[i] = string(rune('a'+i%26)) + "word" + string(rune('0'+i/10)) + string(rune('0'+i%10))
	}

	text := "This is a test text with aword00 and bword11 and some cword22 data for processing."

	kp := NewKeywordProcessor()
	kp.AddKeywordsFromList(keywords).Build()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		kp.ExtractKeywords(text)
	}
}

// 大关键词集 (1000个)
func BenchmarkLargeKeywordSet(b *testing.B) {
	// 生成1000个关键词
	keywords := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		keywords[i] = string(rune('a'+i%26)) + "word" + string(rune('0'+i/100)) + string(rune('0'+(i/10)%10)) + string(rune('0'+i%10))
	}

	text := "This is a test text with aword000 and bword111 and some cword222 data for processing."

	kp := NewKeywordProcessor()
	kp.AddKeywordsFromList(keywords).Build()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		kp.ExtractKeywords(text)
	}
}

// 测试不同文本长度的性能影响

// 短文本 (约50字符)
func BenchmarkShortText(b *testing.B) {
	keywords := []string{"test", "golang", "performance"}
	text := "This is a short test text with golang code."

	kp := NewKeywordProcessor()
	kp.AddKeywordsFromList(keywords).Build()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		kp.ExtractKeywords(text)
	}
}

// 中等文本 (约500字符)
func BenchmarkMediumText(b *testing.B) {
	keywords := []string{"test", "golang", "performance", "data", "processing"}
	text := `This is a medium length test text with multiple sentences.
		It contains golang code and various data processing examples.
		The performance of the keyword extraction should be tested with this text.
		We include multiple occurrences of test keywords in different contexts.
		This helps us understand how the algorithm performs with realistic text lengths.
		The text should be long enough to demonstrate performance characteristics.`

	kp := NewKeywordProcessor()
	kp.AddKeywordsFromList(keywords).Build()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		kp.ExtractKeywords(text)
	}
}

// 长文本 (约5000字符)
func BenchmarkLongText(b *testing.B) {
	keywords := []string{"test", "golang", "performance", "data", "processing"}

	// 重复文本以生成长文本
	baseText := `This is a test sentence with golang performance data processing. `
	text := ""
	for i := 0; i < 80; i++ {
		text += baseText
	}

	kp := NewKeywordProcessor()
	kp.AddKeywordsFromList(keywords).Build()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		kp.ExtractKeywords(text)
	}
}

// 测试大小写敏感 vs 不敏感的性能差异

func BenchmarkCaseInsensitive(b *testing.B) {
	keywords := []string{"Go", "Golang", "Python", "Java", "JavaScript"}
	text := "I love programming in go, golang, PYTHON, java, and javascript."

	kp := NewKeywordProcessor() // 不区分大小写
	kp.AddKeywordsFromList(keywords).Build()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		kp.ExtractKeywords(text)
	}
}

func BenchmarkCaseSensitive(b *testing.B) {
	keywords := []string{"Go", "Golang", "Python", "Java", "JavaScript"}
	text := "I love programming in Go, Golang, Python, Java, and JavaScript."

	kp := NewKeywordProcessor(WithCaseSensitive()) // 区分大小写
	kp.AddKeywordsFromList(keywords).Build()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		kp.ExtractKeywords(text)
	}
}

// 测试中英文混合文本性能
func BenchmarkMixedLanguage(b *testing.B) {
	keywords := []string{"golang", "性能", "测试", "performance", "代码", "code"}
	text := "这是一个golang性能测试的例子。This is a performance test example with 代码 and code mixed together."

	kp := NewKeywordProcessor()
	kp.AddKeywordsFromList(keywords).Build()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		kp.ExtractKeywords(text)
	}
}

// 测试Build性能
func BenchmarkBuild(b *testing.B) {
	keywords := []string{
		"he", "she", "hers", "his", "share", "apple", "banana", "orange",
		"golang", "python", "java", "rust", "c++", "javascript",
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		kp := NewKeywordProcessor()
		kp.AddKeywordsFromList(keywords).Build()
	}
}

// 测试添加关键词的性能
func BenchmarkAddKeywords(b *testing.B) {
	keywords := []string{
		"he", "she", "hers", "his", "share", "apple", "banana", "orange",
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		kp := NewKeywordProcessor()
		for _, word := range keywords {
			kp.AddKeyWord(word)
		}
	}
}
