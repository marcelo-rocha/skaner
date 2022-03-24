package xss_test

import (
	"io"
	"os"
	"path"
	"testing"

	"github.com/marcelo-rocha/skaner/domain/checker"
	"github.com/marcelo-rocha/skaner/domain/checker/xss"
	"github.com/marcelo-rocha/skaner/domain/sourcecode"
	"github.com/stretchr/testify/require"
)

func TestXSSChecker(t *testing.T) {

	testCases := []struct {
		fileName string
		lines    []int
	}{
		{fileName: "xss1.html", lines: []int{13}},
		{fileName: "xss2.html", lines: []int{11}},
		{fileName: "xss3.html", lines: []int{11, 14}},
	}

	workDir, _ := os.Getwd()
	var xssChecker checker.Checker = xss.New()

	for i := range testCases {
		require.True(t, xssChecker.SupportedFileExtension(testCases[i].fileName))
		srcFilePath := path.Join(workDir, "../../../test/data/", testCases[i].fileName)
		f, err := os.Open(srcFilePath)
		require.NoError(t, err)
		defer f.Close()

		content, err := io.ReadAll(f)
		require.NoError(t, err)

		src := sourcecode.NewSourceCode(content, testCases[i].fileName)

		list, err := xssChecker.Check(src)
		require.NoError(t, err)

		require.Len(t, list, len(testCases[i].lines))
		for j := range list {
			require.Equal(t, testCases[i].lines[j], list[j].Line)
		}
	}

}
