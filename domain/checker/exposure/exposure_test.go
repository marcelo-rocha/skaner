package exposure_test

import (
	"io"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/marcelo-rocha/skaner/domain/checker"
	"github.com/marcelo-rocha/skaner/domain/checker/exposure"
	"github.com/marcelo-rocha/skaner/domain/sourcecode"
)

func TestExposure(t *testing.T) {

	testCases := []struct {
		fileName string
		qty      int
		lines    []int
	}{
		{fileName: "exposure1.txt", qty: 1, lines: []int{4}},
		{fileName: "exposure2.txt", qty: 0},
		{fileName: "exposure3.txt", qty: 2, lines: []int{4, 5}},
		{fileName: "exposure4.txt", qty: 1, lines: []int{2}},
	}

	workDir, _ := os.Getwd()
	var exposureChecker checker.Checker
	exposureChecker, err := exposure.New([]string{
		"Bill",
		"$2.2b",
		"Switzerland National Bank",
	})
	require.NoError(t, err)

	for i := range testCases {
		srcFilePath := path.Join(workDir, "../../../test/data/", testCases[i].fileName)
		f, err := os.Open(srcFilePath)
		require.NoError(t, err)
		defer f.Close()

		content, err := io.ReadAll(f)
		require.NoError(t, err)

		src := sourcecode.NewSourceCode(content, testCases[i].fileName)

		list, err := exposureChecker.Check(src)
		require.NoError(t, err)

		require.Len(t, list, testCases[i].qty)
		for j := range list {
			require.Equal(t, list[j].Line, testCases[i].lines[j])
		}
	}

}
