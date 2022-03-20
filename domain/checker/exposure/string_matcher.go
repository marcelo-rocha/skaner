package exposure

import (
	"bytes"

	"github.com/zeebo/xxh3"
)

type MatcherContext interface{}

type Word struct {
	LowerData []byte
	Hash      uint64
}

func ToWord(s []byte) Word {
	lower := bytes.ToLower(s)
	return Word{
		LowerData: lower,
		Hash:      xxh3.Hash(lower),
	}
}

type StringMatcher interface {
	BeginFile() MatcherContext
	NewLine(mc MatcherContext)
	MatchWord(mc MatcherContext, w Word) bool
}
