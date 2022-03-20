package exposure

import "strings"

type StringSequenceMacher struct {
	SimpleMatchers []*SimpleStringMatcher
}

type SequenceContext struct {
	matcherIdx int
}

func NewStringSequenceMatcher(text string) *StringSequenceMacher {
	var matchers []*SimpleStringMatcher
	parts := strings.Split(text, " ")
	for _, s := range parts {
		if s != " " {
			matchers = append(matchers, NewSimpleStringMatcher(s))
		}
	}
	return &StringSequenceMacher{
		SimpleMatchers: matchers,
	}
}

func (sq *StringSequenceMacher) BeginFile() MatcherContext {
	return &SequenceContext{matcherIdx: 0}
}

func (sq *StringSequenceMacher) NewLine(mc MatcherContext) {
	sc := mc.(*SequenceContext)
	sc.matcherIdx = 0
}

func (sq *StringSequenceMacher) MatchWord(mc MatcherContext, w Word) bool {
	sc := mc.(*SequenceContext)

	if sq.SimpleMatchers[sc.matcherIdx].MatchWord(DefaultMatcherContext, w) {
		sc.matcherIdx++
		if sc.matcherIdx == len(sq.SimpleMatchers) {
			return true
		}
	} else {
		sc.matcherIdx = 0
	}
	return false
}
