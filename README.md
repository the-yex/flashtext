# FlashText

<div align="center">

**高性能的 Go 语言 AC 自动机实现**

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.20-blue)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Tests](https://img.shields.io/badge/tests-passing-brightgreen.svg)](.)

*快速、准确、完整的关键词匹配解决方案*

[特性](#特性) • [安装](#安装) • [快速开始](#快速开始) • [性能](#性能) • [文档](#文档)

</div>

---

## 📖 简介

FlashText 是一个高性能的 **Aho-Corasick 自动机**实现，专为 Go 语言设计。使用经典的 AC 自动机算法，能够在文本中快速查找和提取大量关键词，并保证找到**所有重叠匹配**。

### 为什么选择 FlashText？

- ⚡ **性能卓越**: 比正则表达式快 **60倍**
- 🎯 **完整匹配**: 找到所有重叠匹配，不遗漏任何结果
- 📦 **简单易用**: 清晰的 API，5 行代码即可上手
- 🔧 **生产就绪**: 经过完整测试，可用于敏感词过滤、内容审核等场景
- 🌏 **多语言支持**: 完整的 UTF-8 和中文支持

---

## ✨ 特性

### 核心功能

- ✅ **完整的 AC 自动机实现** - 包含 Trie 树和失败指针机制
- ✅ **🧠 自适应容量引擎** - 越用越快，智能预测内存分配
- ✅ **重叠匹配检测** - 找到所有可能的匹配，包括重叠部分
- ✅ **大小写控制** - 支持大小写敏感/不敏感模式
- ✅ **高效批量处理** - 一次添加多个关键词
- ✅ **UTF-8 完整支持** - 正确处理中文、日韩文等多字节字符

### 性能特点


| 场景                   | 性能       |
| ---------------------- | ---------- |
| 1000个关键词 + 6MB文本 | ~380ms     |
| vs 正则表达式          | **快60倍** |
| 时间复杂度             | O(n) 线性  |

📊 详细性能数据请参考 [PERFORMANCE.md](PERFORMANCE.md)

---

## 🚀 安装

```bash
go get github.com/the-yex/flashtext
```

**要求**: Go 1.20+

---

## 💡 快速开始

### 基础用法

```go
package main

import (
    "fmt"
    "github.com/the-yex/flashtext"
)

func main() {
    // 1. 创建处理器（不区分大小写）
    kp := flashtext.NewKeywordProcessor(false)
  
    // 2. 添加关键词并构建
    kp.AddKeywordsFromList([]string{"golang", "python", "java"}).Build()
  
    // 3. 提取关键词
    text := "I love Golang and Python programming!"
    matches := kp.ExtractKeywords(text)
  
    // 4. 处理结果
    for _, match := range matches {
        fmt.Printf("找到: %s [%d:%d]\n", 
            match.MatchString(), match.Start(), match.End())
    }
}
```

**输出**:

```
找到: Golang [7:13]
找到: Python [18:24]
```

### 高级用法

#### 大小写敏感匹配

```go
kp := flashtext.NewKeywordProcessor(true) // 区分大小写
kp.AddKeyWord("Go").Build()

kp.ExtractKeywords("I use Go and go") 
// 只匹配 "Go"，不匹配 "go"
```

#### 处理字节数组

```go
data := []byte("some binary data with keywords")
matches := kp.ExtractKeywordsFromBytes(data)
```

#### 链式调用

```go
kp := flashtext.NewKeywordProcessor(false).
    AddKeyWord("apple").
    AddKeyWord("banana").
    AddKeywordsFromList([]string{"orange", "grape"}).
    Build()
```

---

## 🎯 使用场景

### ✅ 推荐场景

- **敏感词过滤** - 内容审核、评论检测
- **文本分析** - 关键词提取、实体识别
- **数据挖掘** - 大规模文本处理
- **SEO工具** - 关键词密度分析
- **日志分析** - 错误关键词检测

### ⚠️ 不适合的场景

- 关键词数量 < 10 (正则表达式可能更简单)
- 只需要部分匹配 (可考虑简化的 Trie 实现)

---

## 📊 性能

### 性能对比 (1000个关键词)


| 实现          | 时间      | 内存  | 完整性    |
| ------------- | --------- | ----- | --------- |
| **FlashText** | **383ms** | 172MB | ✅ 完整   |
| 正则表达式    | 22,900ms  | 1MB   | ✅ 完整   |
| 简化Trie      | 172ms     | 133MB | ❌ 不完整 |

### 重叠匹配示例

在文本 `"hershey"` 中查找 `["he", "she", "hers"]`:

```
FlashText (AC自动机):
✓ he    [0:2]
✓ hers  [0:4]  
✓ she   [3:6]   ← 重叠匹配
✓ he    [4:6]   ← 重叠匹配
总计: 4个

简化Trie实现:
✓ he    [0:1]
✓ hers  [0:3]
✓ he    [4:5]
❌ 遗漏 "she"
总计: 3个
```

**结论**: FlashText 保证找到所有匹配，这在敏感词过滤等安全场景下至关重要。

📖 更多性能数据和分析，请查看 [PERFORMANCE.md](PERFORMANCE.md)

### 🧠 自适应容量引擎

本库引入了创新的 **EWMA (指数加权移动平均)** 算法来动态优化内存分配。

- **智能学习**: 自动学习您的业务数据的关键词密度。
- **越用越快**: 随着运行时间增长，内存预分配越来越精准，扩容次数趋近于 0。
- **无感运行**: 全自动后台优化，无需任何配置。

**实测数据 (1000次调用)**:
容量在**前50次**调用中迅速学习并稳定，后续**950次**调用中扩容次数为**0**。

| 迭代次数 | 预估容量 | 实际匹配 | 状态 |
| :--- | :--- | :--- | :--- |
| 1 | 215 | 200 | 学习中 |
| 10 | 383 | 200 | 快速适应 |
| **50+** | **423** | 200 | **完全稳定** |

👉 了解更多设计细节: [ADAPTIVE_ENGINE.md](ADAPTIVE_ENGINE.md)

### ⚠️ 资源释放 (重要)

由于自适应引擎在后台运行一个轻量级统计协程，**如果您的服务在运行过程中需要频繁销毁或重新构建 `KeywordProcessor`**，请务必调用 `Close()` 方法以释放资源：

```go
kp := flashtext.NewKeywordProcessor()
// ... 使用 ...

// 当不再需要该实例，或者即将创建新实例替换它时：
kp.Close() 
```

对于作为**全局单例**长期运行的 `KeywordProcessor`，则无需时刻关注 `Close`。

---

## 📚 文档

### API 文档

#### 创建处理器

```go
// 参数: caseSensitive - 是否区分大小写
kp := flashtext.NewKeywordProcessor(caseSensitive bool)
```

#### 添加关键词

```go
// 添加单个关键词
kp.AddKeyWord(keyword string) *KeywordProcessor

// 批量添加关键词
kp.AddKeywordsFromList(keywords []string) *KeywordProcessor
```

#### 构建索引

```go
// 必须在添加完所有关键词后调用
kp.Build()
```

#### 提取关键词

```go
// 从字符串提取
matches := kp.ExtractKeywords(text string) []Match

// 从字节数组提取
matches := kp.ExtractKeywordsFromBytes(data []byte) []Match
```

#### Match 结构

```go
type Match struct {
    match string  // 匹配的文本
    start int     // 开始位置（字节）
    end   int     // 结束位置（字节）
}

// 获取方法
match.MatchString() string  // 匹配的文本
match.Start() int           // 开始位置
match.End() int             // 结束位置
```

### 完整示例

请参考测试文件 `keywordprocessor_test.go`

---

## 🧪 测试

### 运行测试

```bash
# 功能测试
go test -v

# 性能测试
go test -bench=. -benchmem

# 性能对比测试
go test -bench=BenchmarkComparison -benchmem -run=^$ -timeout=30m
```

### 测试覆盖

- ✅ 基础匹配测试
- ✅ 中文字符测试
- ✅ 大小写敏感/不敏感测试
- ✅ 边缘情况测试
- ✅ 重叠匹配测试
- ✅ 性能对比测试

---

## 🔄 与其他实现的对比

### vs Python FlashText

Python 版本在连续文本中存在匹配 bug（需要分词），而本实现能正确处理连续文本：

```python
# Python FlashText
kp.extract_keywords('ahishehersshare')  # 返回 []  ❌

kp.extract_keywords('a his he hers share')  # 返回 ['his', 'he', 'hers', 'share']  ✅
```

```go
// FlashText Go
kp.ExtractKeywords("ahishershare")  // 正确匹配所有  ✅
```

### vs ayoyu/flashtext (Go)

ayoyu/flashtext 使用简化的 Trie 树，速度更快但会遗漏重叠匹配。FlashText 使用完整的 AC 自动机，保证匹配完整性。

**选择建议**:

- 需要完整匹配 → FlashText
- 只追求速度且不关心重叠 → ayoyu/flashtext

---

## 🛠️ 最佳实践

### 1. 复用实例

```go
// ✅ 好的做法：复用实例
kp := flashtext.NewKeywordProcessor(false)
kp.AddKeywordsFromList(keywords).Build()

for _, text := range texts {
    matches := kp.ExtractKeywords(text)
}

// ❌ 不好的做法：每次都创建新实例
for _, text := range texts {
    kp := flashtext.NewKeywordProcessor(false)
    kp.AddKeywordsFromList(keywords).Build()
    matches := kp.ExtractKeywords(text)
}
```

### 2. 大文本分块处理

```go
const chunkSize = 1024 * 1024 // 1MB
for i := 0; i < len(text); i += chunkSize {
    end := min(i+chunkSize, len(text))
    matches := kp.ExtractKeywords(text[i:end])
    // 处理 matches...
}
```

---

## 📖 设计理念

### 为什么不支持关键词替换？

我们遵循**单一职责原则**，专注于做好一件事：**高效、准确的关键词匹配**。

关键词替换可以很容易地在应用层实现：

```go
matches := kp.ExtractKeywords(text)
// 用 matches 信息自行替换
```

这种设计让库保持简洁，同时给用户最大的灵活性。

---

## 🤝 贡献

我们欢迎各种形式的贡献！

### 如何贡献

1. Fork 本仓库
2. 创建您的特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交您的更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启一个 Pull Request

### 贡献指南

- ✅ 添加测试用例
- ✅ 保持代码风格一致
- ✅ 更新相关文档
- ✅ 确保所有测试通过

---

## 📄 许可证

本项目采用 [MIT 许可证](LICENSE)

---

## 🙏 致谢

- 感谢 [Aho-Corasick](https://en.wikipedia.org/wiki/Aho%E2%80%93Corasick_algorithm) 算法的发明者
- 参考了 Python flashtext 库的设计思路

---

## 📮 联系方式

- **Issues**: [GitHub Issues](https://github.com/the-yex/flashtext/issues)
- **Discussions**: [GitHub Discussions](https://github.com/the-yex/flashtext/discussions)

---

<div align="center">

**如果这个项目对您有帮助，请给我们一个 ⭐ Star！**

Made with ❤️ by the FlashText team

</div>
