package benchmark_test

import (
	"fmt"
	"testing"

	ayoyu "github.com/ayoyu/flashtext"
	"github.com/the-yex/flashtext"
)

// 验证ayoyu是否会漏掉重叠匹配
func TestAyoyuOverlapMatching(t *testing.T) {
	// 测试重叠匹配：在"hershey"中应该匹配["he", "she", "hers", "he"]
	keywords := []string{"he", "she", "hers"}
	text := "hershey"

	// 本库测试
	kp := flashtext.NewKeywordProcessor()
	kp.AddKeywordsFromList(keywords).Build()
	ourMatches := kp.ExtractKeywords(text)

	fmt.Printf("\n本库匹配结果 (AC自动机):\n")
	for _, m := range ourMatches {
		fmt.Printf("  - %s [%d:%d]\n", m.MatchString(), m.Start(), m.End())
	}

	// ayoyu测试
	ayoyuFlash := ayoyu.NewFlashKeywords(false)
	for _, kw := range keywords {
		ayoyuFlash.Add(kw)
	}
	ayoyuMatches := ayoyuFlash.Search(text)

	fmt.Printf("\nayoyu匹配结果 (Trie树):\n")
	for _, m := range ayoyuMatches {
		fmt.Printf("  - %s [%d:%d]\n", m.Key, m.Start, m.End)
	}

	fmt.Printf("\n结果对比:\n")
	fmt.Printf("  本库: %d个匹配\n", len(ourMatches))
	fmt.Printf("  ayoyu: %d个匹配\n", len(ayoyuMatches))

	if len(ourMatches) != len(ayoyuMatches) {
		fmt.Printf("  ⚠️  匹配数量不同！本库找到%d个，ayoyu找到%d个\n", len(ourMatches), len(ayoyuMatches))
	}
}
