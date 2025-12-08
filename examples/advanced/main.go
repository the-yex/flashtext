package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/the-yex/flashtext"
)

func main() {
	// 演示：重叠匹配和批量处理
	// Demo: Overlapping matches and batch processing

	fmt.Println("=== 1. Overlapping Matches / 重叠匹配 ===")
	kp := flashtext.NewKeywordProcessor()
	
	// 添加一组包含包含关系的关键词
	keywords := []string{"sys", "system", "tem", "operating system"}
	kp.AddKeywordsFromList(keywords).Build()

	text := "The operating system is complex."
	fmt.Println("Keywords:", keywords)
	fmt.Println("Text:", text)
	
	matches := kp.ExtractKeywords(text)
	for _, m := range matches {
		fmt.Printf("Found: %s \t[%d:%d]\n", m.MatchString(), m.Start(), m.End())
	}
	// 预期结果：会找到 "operating system", "system", "sys", "tem" 等所有匹配
	// Expected: finds all overlapping keywords

	fmt.Println("\n=== 2. Processing Bytes / 处理字节数据 ===")
	// 模拟读取文件内容
	filename := "test_sample.txt"
	content := []byte("Hello world, this is a binary file test.")
	_ = ioutil.WriteFile(filename, content, 0644)
	defer os.Remove(filename)

	fileData, _ := ioutil.ReadFile(filename)
	
	kp2 := flashtext.NewKeywordProcessor()
	kp2.AddKeywordsFromList([]string{"hello", "file", "test"}).Build()
	
	// 直接处理 []byte，避免转换string的开销
	byteMatches := kp2.ExtractKeywordsFromBytes(fileData)
	for _, m := range byteMatches {
		fmt.Printf("Found in bytes: %s\n", m.MatchString())
	}
}
