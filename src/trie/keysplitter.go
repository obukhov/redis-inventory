package trie

// Splitter abstraction to split string in fragments
type Splitter interface {
	// Split splits string key to fragments with given strategy
	Split(in string) []string
}

// PunctuationSplitter splitting keys by a specific set of symbols (i.e. punctuation)
type PunctuationSplitter struct {
	dividers map[rune]bool
}

// NewPunctuationSplitter creates PunctuationSplitter
func NewPunctuationSplitter(punctuation ...rune) *PunctuationSplitter {
	m := make(map[rune]bool)
	for _, r := range punctuation {
		m[r] = true
	}

	return &PunctuationSplitter{dividers: m}
}

// Split splits string key to fragments with given strategy
func (s *PunctuationSplitter) Split(in string) []string {
	result := make([]string, 0)

	cur := ""
	isPunkt := false

	for _, c := range in {
		if _, found := s.dividers[c]; found {
			isPunkt = true
			cur = cur + string(c)
		} else {
			if isPunkt {
				result = append(result, cur)
				cur = string(c)
				isPunkt = false
			} else {
				cur = cur + string(c)
			}
		}
	}

	result = append(result, cur)

	return result
}
