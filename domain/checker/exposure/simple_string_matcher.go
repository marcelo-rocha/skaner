package exposure

import (
	"bytes"

	"github.com/zeebo/xxh3"
)

var DefaultMatcherContext MatcherContext = nil

type SimpleStringMatcher struct {
	lowerData []byte
	hash      uint64
}

func NewSimpleStringMatcher(s string) *SimpleStringMatcher {
	lowerData := bytes.ToLower([]byte(s))
	return &SimpleStringMatcher{
		lowerData: lowerData,
		hash:      xxh3.Hash(lowerData),
	}
}

func (m *SimpleStringMatcher) BeginFile() MatcherContext {
	return DefaultMatcherContext
}

func (m *SimpleStringMatcher) NewLine(mc MatcherContext) {
}

func (m *SimpleStringMatcher) MatchWord(mc MatcherContext, w Word) bool {
	return m.hash == w.Hash && bytes.Equal(m.lowerData, w.LowerData)
}
