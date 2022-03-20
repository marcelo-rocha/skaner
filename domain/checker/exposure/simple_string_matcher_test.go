package exposure_test

import (
	"testing"

	"github.com/marcelo-rocha/skaner/domain/checker/exposure"
	"github.com/stretchr/testify/require"
)

func TestSimpleMatch(t *testing.T) {
	matcher := exposure.NewSimpleStringMatcher("Confidential")

	mc := matcher.BeginFile()
	w := exposure.ToWord([]byte("Secret"))
	require.False(t, matcher.MatchWord(mc, w))

	w = exposure.ToWord([]byte("Confidential"))
	require.True(t, matcher.MatchWord(mc, w))

	w = exposure.ToWord([]byte("confidential"))
	require.True(t, matcher.MatchWord(mc, w))

	w = exposure.ToWord([]byte("CONFIDENTIAL"))
	require.True(t, matcher.MatchWord(mc, w))
}
