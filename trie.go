package flashtext

/*
* @Author: zouyx
* @Email:
* @Date:   2023/5/17 17:27
* @Package:
 */

type Node struct {
	children map[rune]*Node // 使用 map 存储叶子节点,key:'char' ,value: *Node
	exist    []int          // 关键词节点是一个完整的匹配词，记录其长度  可以在build 的时候去重，匹配的时候遍历比map快
	failure  *Node          // 记录失败指针
}

func newNode() *Node {
	return &Node{
		children: make(map[rune]*Node),
		exist:    nil,
	}
}

type Match struct {
	match string
	start int
	end   int
}

func (m *Match) MatchString() string {
	return m.match
}

func (m *Match) Start() int {
	return m.start
}

func (m *Match) End() int {
	return m.end
}
