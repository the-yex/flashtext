package flashtext

import (
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
	root          *Node
	caseSensitive bool // 匹配是否区分大小写
}

// NewKeywordProcessor creates a new processor instance.
// caseSensitive: if true, matches are case-sensitive.
func NewKeywordProcessor(caseSensitive bool) *KeywordProcessor {
	return &KeywordProcessor{
		root:          newNode(),
		caseSensitive: caseSensitive,
	}
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
	node.exist[len(keyword)] = struct{}{}
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
			for key := range childNode.failure.exist {
				childNode.exist[key] = struct{}{}
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
func (kp *KeywordProcessor) walk(sentence string, wf WalkFn) {
	currentNode := kp.root
	idx := 0
	for len(sentence) > 0 {
		r, size := utf8.DecodeRuneInString(sentence)
		idx += size
		sentence = sentence[size:]

		if !kp.caseSensitive {
			r = unicode.ToLower(r)
		}

		// 失败指针回溯
		for currentNode.children[r] == nil && currentNode.failure != nil {
			currentNode = currentNode.failure
		}

		if currentNode.children[r] == nil {
			continue
		}

		currentNode = currentNode.children[r]

		// 输出所有匹配
		for length := range currentNode.exist {
			if !wf(idx-length, idx) {
				return
			}
		}
	}
}

// ExtractKeywords searches for keywords in a string.
// It returns a slice of all matches found.
func (kp *KeywordProcessor) ExtractKeywords(sentence string) []Match {
	// 优化: 预分配容量
	matches := make([]Match, 0, 16)
	if len(sentence) == 0 {
		return matches
	}

	kp.walk(sentence, func(start, end int) bool {
		matches = append(matches, Match{
			start: start,
			end:   end,
			match: sentence[start:end],
		})
		return true
	})
	return matches
}

// ExtractKeywordsFromBytes searches for keywords in a byte slice.
// It returns a slice of all matches found.
func (kp *KeywordProcessor) ExtractKeywordsFromBytes(sentence []byte) []Match {
	// 优化: 预分配容量 + 统一使用 walk
	matches := make([]Match, 0, 16)
	if len(sentence) == 0 {
		return matches
	}

	str := string(sentence)
	kp.walk(str, func(start, end int) bool {
		matches = append(matches, Match{
			start: start,
			end:   end,
			match: str[start:end],
		})
		return true
	})
	return matches
}
