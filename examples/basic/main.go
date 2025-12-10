package main

import (
	"fmt"

	"github.com/the-yex/flashtext"
)

func main() {
	// 1. 初始化处理器 (true表示区分大小写)
	// Initialize processor (true for case sensitive)
	kp := flashtext.NewKeywordProcessor(flashtext.WithCaseSensitive())
	defer kp.Close()
	// 2. 添加关键词
	// Add keywords
	kp.AddKeyWord("Go")
	kp.AddKeyWord("Python")
	kp.AddKeywordsFromList([]string{"Java", "C++", "Rust"})

	// 3. 构建索引 (必须调用!)
	// Build the index (Required!)
	kp.Build()

	// 4. 准备文本
	// Prepare text
	text := "I am learning Go and Python. python is also good, but I like Rust more."

	// 5. 提取关键词
	// Extract keywords
	fmt.Println("Text:", text)
	fmt.Println("--- Found Keywords ---")
	matches := kp.ExtractKeywords(text)

	for _, match := range matches {
		fmt.Printf("Keyword: %-10s at [%d:%d]\n", match.MatchString(), match.Start(), match.End())
	}
}
