package exposure_test

import (
	"testing"

	"github.com/marcelo-rocha/skaner/domain/checker/exposure"
	"github.com/stretchr/testify/require"
)

func TestSequenceMatcher(t *testing.T) {
	var sm exposure.StringMatcher = exposure.NewStringSequenceMatcher("ACME & Associates")
	mc := sm.BeginFile()

	require.False(t, sm.MatchWord(mc, exposure.ToWord([]byte("News"))))
	require.False(t, sm.MatchWord(mc, exposure.ToWord([]byte("ACME"))))
	require.False(t, sm.MatchWord(mc, exposure.ToWord([]byte("&"))))
	require.True(t, sm.MatchWord(mc, exposure.ToWord([]byte("associates"))))
}

func TestSequenceMatcherMultiline(t *testing.T) {
	var sm exposure.StringMatcher = exposure.NewStringSequenceMatcher("ACME & Associates")
	mc := sm.BeginFile()

	require.False(t, sm.MatchWord(mc, exposure.ToWord([]byte("News"))))
	require.False(t, sm.MatchWord(mc, exposure.ToWord([]byte("about"))))
	require.False(t, sm.MatchWord(mc, exposure.ToWord([]byte("ACME"))))
	require.False(t, sm.MatchWord(mc, exposure.ToWord([]byte("&"))))
	sm.NewLine(mc)
	require.False(t, sm.MatchWord(mc, exposure.ToWord([]byte("associates"))))
	require.False(t, sm.MatchWord(mc, exposure.ToWord([]byte("acme"))))
	require.False(t, sm.MatchWord(mc, exposure.ToWord([]byte("&"))))
	require.True(t, sm.MatchWord(mc, exposure.ToWord([]byte("associates"))))
}
