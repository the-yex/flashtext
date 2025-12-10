package flashtext

import (
	"context"
	"unicode"
	"unicode/utf8"
)

// Package flashtext implements the Aho-Corasick algorithm for efficient keyword matching.
// It is optimized for finding multiple patterns in a text simultaneously.

// WalkFn is the callback function used during traversal.
// It receives the start and end byte positions of the match.
// Return false to stop traversal.
type WalkFn func(start, end int) bool

// KeywordProcessor controls the keyword matching process.
// It holds the AC automaton trie and configuration.
type KeywordProcessor struct {
	cancel        context.CancelFunc
	root          *Node
	stats         *stats // 异步统计模块，根据词库动态调整 density ，跑的越久性能越好
	caseSensitive bool   // 匹配是否区分大小写
	matchDensity  float64
}
type Option func(*KeywordProcessor)

// 是否区分大小写，默认不区分大小写
func WithCaseSensitive() Option {
	return func(processor *KeywordProcessor) {
		processor.caseSensitive = true
	}
}

// NewKeywordProcessor creates a new processor instance.
// caseSensitive: if true, matches are case-sensitive.
func NewKeywordProcessor(opts ...Option) *KeywordProcessor {
	ctx := context.Background()
	processor := &KeywordProcessor{
		root:          newNode(),
		caseSensitive: false,
		stats:         newStats(ctx, defaultAlpha, defaultBuffer),
	}
	for _, opt := range opts {
		opt(processor)
	}
	return processor
}

func (kp *KeywordProcessor) setItem(keyword string) {
	if len(keyword) == 0 {
		return
	}

	node := kp.root
	for _, char := range keyword {
		if !kp.caseSensitive {
			char = unicode.ToLower(char)
		}
		if _, ok := node.children[char]; !ok {
			node.children[char] = newNode()
		}
		node = node.children[char]
	}
	// 记录当前匹配词的长度
	node.exist = append(node.exist, len([]rune(keyword)))
}

// Build constructs the failure pointers for the AC automaton.
// This MUST be called after all keywords are added and before matching.
func (kp *KeywordProcessor) Build() {
	// 优化: 预分配队列容量
	queue := make([]*Node, 0, 128)
	queue = append(queue, kp.root)

	for len(queue) > 0 {
		currentNode := queue[0]
		queue = queue[1:]

		for char, childNode := range currentNode.children {
			queue = append(queue, childNode)
			faFail := currentNode.failure

			for faFail != nil && faFail.children[char] == nil {
				faFail = faFail.failure
			}
			childNode.failure = kp.root
			if faFail != nil {
				childNode.failure = faFail.children[char]
			}
			// Merge exist and deduplicate
			childNode.exist = append(childNode.exist, childNode.failure.exist...)
			tmp := make(map[int]struct{}, len(childNode.exist))
			for _, l := range childNode.exist {
				tmp[l] = struct{}{}
			}
			if len(tmp) < len(childNode.exist) {
				childNode.exist = childNode.exist[:0]
				for l := range tmp {
					childNode.exist = append(childNode.exist, l)
				}
			}
		}
	}
}

// AddKeyWord adds a single keyword to the processor.
// Returns the processor for chaining.
func (kp *KeywordProcessor) AddKeyWord(keyword string) *KeywordProcessor {
	kp.setItem(keyword)
	return kp
}

// AddKeywordsFromList adds multiple keywords from a slice.
// Returns the processor for chaining.
func (kp *KeywordProcessor) AddKeywordsFromList(keywords []string) *KeywordProcessor {
	for _, keyword := range keywords {
		kp.setItem(keyword)
	}
	return kp
}

// walk
func (kp *KeywordProcessor) walk(sentence []rune, wf WalkFn) {
	node := kp.root

	for i, r := range sentence {
		if !kp.caseSensitive {
			r = unicode.ToLower(r)
		}
		for node.children[r] == nil && node != kp.root {
			node = node.failure
		}

		if node.children[r] != nil {
			node = node.children[r]
		}

		for _, l := range node.exist {
			if !wf(i+1-l, i+1) {
				return
			}
		}
	}
}

// ExtractKeywords searches for keywords in a string.
// It returns a slice of all matches found.
func (kp *KeywordProcessor) ExtractKeywords(sentence string) []Match {
	// 优化: 预分配容量
	runes := []rune(sentence)
	if len(runes) == 0 {
		return nil
	}
	density := kp.stats.getDensity()
	capEstimate := int(float64(len(runes)) * density)
	matches := make([]Match, 0, capEstimate+1)
	byteOffsets := make([]int, len(runes)+1)
	for i, r := range runes {
		byteOffsets[i+1] = byteOffsets[i] + utf8.RuneLen(r)
	}
	kp.walk(runes, func(start, end int) bool {
		startByte := byteOffsets[start]
		endByte := byteOffsets[end]
		matches = append(matches, Match{
			start: startByte,
			end:   endByte,
			match: sentence[startByte:endByte],
		})
		return true
	})
	kp.stats.add(len(matches), len(runes))
	return matches
}

// ExtractKeywordsFromBytes searches for keywords in a byte slice.
// It returns a slice of all matches found.
func (kp *KeywordProcessor) ExtractKeywordsFromBytes(sentence []byte) []Match {
	// 优化: 预分配容量 + 统一使用 walk
	return kp.ExtractKeywords(string(sentence))
}

func (kp *KeywordProcessor) Close() {
	if kp.stats != nil {
		kp.stats.close()
	}
}
