package xss_test

import (
	"testing"

	"github.com/marcelo-rocha/skaner/domain/checker/xss"
	"github.com/stretchr/testify/require"
)

func TestScanJS(t *testing.T) {
	src := `
	// generating  a random number
	const a = Math.random();
	console.log(a);	
	// calculate other random
	console.log(Math.random());	
	`
	lines, err := xss.ScanJS([]byte(src), []byte("random"))
	require.NoError(t, err)
	require.Equal(t, 2, len(lines))
	require.Equal(t, 3, lines[0])
	require.Equal(t, 6, lines[1])
}
