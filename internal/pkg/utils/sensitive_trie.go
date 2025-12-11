package utils

// SensitiveTrie 敏感词前缀树
type SensitiveTrie struct {
	root *TrieNode
}

type TrieNode struct {
	children map[rune]*TrieNode
	isEnd    bool
}

func NewSensitiveTrie() *SensitiveTrie {
	return &SensitiveTrie{
		root: &TrieNode{
			children: make(map[rune]*TrieNode),
			isEnd:    false,
		},
	}
}

// AddWord 添加敏感词
func (t *SensitiveTrie) AddWord(word string) {
	node := t.root
	for _, char := range word {
		if _, ok := node.children[char]; !ok {
			node.children[char] = &TrieNode{
				children: make(map[rune]*TrieNode),
				isEnd:    false,
			}
		}
		node = node.children[char]
	}
	node.isEnd = true
}

// Validate 验证文本，返回包含的敏感词
func (t *SensitiveTrie) Validate(text string) []string {
	var foundWords []string
	chars := []rune(text)
	length := len(chars)

	for i := 0; i < length; i++ {
		node := t.root
		for j := i; j < length; j++ {
			char := chars[j]
			if child, ok := node.children[char]; ok {
				node = child
				if node.isEnd {
					foundWords = append(foundWords, string(chars[i:j+1]))
					// Optionally continue to find longer words or stop?
					// Usually we want all matches or longest?
					// RuoYi: Default returns all.
				}
			} else {
				break
			}
		}
	}
	// Deduplicate
	return uniqueStrings(foundWords)
}

func uniqueStrings(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
