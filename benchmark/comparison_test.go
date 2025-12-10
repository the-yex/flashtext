package benchmark_test

import (
	"bufio"
	"os"
	"regexp"
	"strings"
	"testing"

	ayoyu "github.com/ayoyu/flashtext"
	"github.com/the-yex/flashtext"
)

// 加载测试数据
func loadWords(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if word != "" {
			words = append(words, word)
		}
	}
	return words, scanner.Err()
}

func loadCorpus(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ============================================
// 性能对比测试：本库 vs Regex vs Ayoyu/flashtext
// ============================================

// 使用真实测试数据进行对比
func BenchmarkComparison_ThisLibrary(b *testing.B) {
	words, err := loadWords("./benchmark_data/words_benchmark_test.txt")
	if err != nil {
		b.Skip("测试数据文件不存在:", err)
		return
	}
	corpus, err := loadCorpus("./benchmark_data/corpus_benchmark_test.txt")
	if err != nil {
		b.Skip("语料文件不存在:", err)
		return
	}

	// 只使用前10000个关键词进行测试（避免太慢）
	if len(words) > 10000 {
		words = words[:10000]
	}

	kp := flashtext.NewKeywordProcessor()
	kp.AddKeywordsFromList(words).Build()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = kp.ExtractKeywords(corpus)
	}
}

func BenchmarkComparison_Regex(b *testing.B) {
	words, err := loadWords("./benchmark_data/words_benchmark_test.txt")
	if err != nil {
		b.Skip("测试数据文件不存在:", err)
		return
	}
	corpus, err := loadCorpus("./benchmark_data/corpus_benchmark_test.txt")
	if err != nil {
		b.Skip("语料文件不存在:", err)
		return
	}

	// 只使用前10000个关键词
	if len(words) > 10000 {
		words = words[:10000]
	}

	// 构建正则表达式: \b(word1|word2|word3)\b
	pattern := "\\b(" + strings.Join(words, "|") + ")\\b"
	re := regexp.MustCompile(pattern)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = re.FindAllString(corpus, -1)
	}
}

func BenchmarkComparison_AyoyuFlashtext(b *testing.B) {
	words, err := loadWords("./benchmark_data/words_benchmark_test.txt")
	if err != nil {
		b.Skip("测试数据文件不存在:", err)
		return
	}
	corpus, err := loadCorpus("./benchmark_data/corpus_benchmark_test.txt")
	if err != nil {
		b.Skip("语料文件不存在:", err)
		return
	}

	// 只使用前10000个关键词
	if len(words) > 10000 {
		words = words[:10000]
	}

	flash := ayoyu.NewFlashKeywords(false)
	for _, word := range words {
		flash.Add(word)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = flash.Search(corpus)
	}
}

// ============================================
// 小规模对比测试（100个关键词）
// ============================================

func BenchmarkSmallScale_ThisLibrary(b *testing.B) {
	words, err := loadWords("./benchmark_data/words_benchmark_test.txt")
	if err != nil {
		b.Skip("测试数据文件不存在:", err)
		return
	}
	corpus, err := loadCorpus("./benchmark_data/corpus_benchmark_test.txt")
	if err != nil {
		b.Skip("语料文件不存在:", err)
		return
	}

	// 只使用100个关键词
	if len(words) > 100 {
		words = words[:100]
	}

	kp := flashtext.NewKeywordProcessor()
	kp.AddKeywordsFromList(words).Build()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = kp.ExtractKeywords(corpus)
	}
}

func BenchmarkSmallScale_Regex(b *testing.B) {
	words, err := loadWords("./benchmark_data/words_benchmark_test.txt")
	if err != nil {
		b.Skip("测试数据文件不存在:", err)
		return
	}
	corpus, err := loadCorpus("./benchmark_data/corpus_benchmark_test.txt")
	if err != nil {
		b.Skip("语料文件不存在:", err)
		return
	}

	if len(words) > 100 {
		words = words[:100]
	}

	pattern := "\\b(" + strings.Join(words, "|") + ")\\b"
	re := regexp.MustCompile(pattern)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = re.FindAllString(corpus, -1)
	}
}

func BenchmarkSmallScale_AyoyuFlashtext(b *testing.B) {
	words, err := loadWords("./benchmark_data/words_benchmark_test.txt")
	if err != nil {
		b.Skip("测试数据文件不存在:", err)
		return
	}
	corpus, err := loadCorpus("./benchmark_data/corpus_benchmark_test.txt")
	if err != nil {
		b.Skip("语料文件不存在:", err)
		return
	}

	if len(words) > 100 {
		words = words[:100]
	}

	flash := ayoyu.NewFlashKeywords(false)
	for _, word := range words {
		flash.Add(word)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = flash.Search(corpus)
	}
}

// ============================================
// 中规模对比测试（1000个关键词）
// ============================================

func BenchmarkMediumScale_ThisLibrary(b *testing.B) {
	words, err := loadWords("./benchmark_data/words_benchmark_test.txt")
	if err != nil {
		b.Skip("测试数据文件不存在:", err)
		return
	}
	corpus, err := loadCorpus("./benchmark_data/corpus_benchmark_test.txt")
	if err != nil {
		b.Skip("语料文件不存在:", err)
		return
	}

	if len(words) > 1000 {
		words = words[:1000]
	}

	kp := flashtext.NewKeywordProcessor()
	kp.AddKeywordsFromList(words).Build()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = kp.ExtractKeywords(corpus)
	}
}

func BenchmarkMediumScale_Regex(b *testing.B) {
	words, err := loadWords("./benchmark_data/words_benchmark_test.txt")
	if err != nil {
		b.Skip("测试数据文件不存在:", err)
		return
	}
	corpus, err := loadCorpus("./benchmark_data/corpus_benchmark_test.txt")
	if err != nil {
		b.Skip("语料文件不存在:", err)
		return
	}

	if len(words) > 1000 {
		words = words[:1000]
	}

	pattern := "\\b(" + strings.Join(words, "|") + ")\\b"
	re := regexp.MustCompile(pattern)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = re.FindAllString(corpus, -1)
	}
}

func BenchmarkMediumScale_AyoyuFlashtext(b *testing.B) {
	words, err := loadWords("./benchmark_data/words_benchmark_test.txt")
	if err != nil {
		b.Skip("测试数据文件不存在:", err)
		return
	}
	corpus, err := loadCorpus("./benchmark_data/corpus_benchmark_test.txt")
	if err != nil {
		b.Skip("语料文件不存在:", err)
		return
	}

	if len(words) > 1000 {
		words = words[:1000]
	}

	flash := ayoyu.NewFlashKeywords(false)
	for _, word := range words {
		flash.Add(word)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = flash.Search(corpus)
	}
}
