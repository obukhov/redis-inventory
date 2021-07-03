package trie

type Splitter interface {
	Split(in string) []string
}

type PunctuationSplitter struct {
	dividers map[rune]bool
}

func NewPunctuationSplitter(punctuation []rune) *PunctuationSplitter {
	m := make(map[rune]bool)
	for _, r := range punctuation {
		m[r] = true
	}

	return &PunctuationSplitter{dividers: m}
}

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
