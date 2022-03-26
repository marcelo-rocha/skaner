package sql_test

import (
	"io"
	"os"
	"path"
	"testing"

	"github.com/marcelo-rocha/skaner/domain/checker"
	"github.com/marcelo-rocha/skaner/domain/checker/sql"
	"github.com/marcelo-rocha/skaner/domain/sourcecode"
	"github.com/stretchr/testify/require"
)

func TestSQL(t *testing.T) {
	testCases := []struct {
		fileName string
		lines    []int
	}{
		{fileName: "sql1.go", lines: []int{7}},
		{fileName: "sql2.go", lines: []int{7, 15}},
	}

	workDir, _ := os.Getwd()
	var sqlChecker checker.Checker = sql.New()

	for i := range testCases {
		require.True(t, sqlChecker.SupportedFileExtension(testCases[i].fileName))
		srcFilePath := path.Join(workDir, "../../../test/data/", testCases[i].fileName)
		f, err := os.Open(srcFilePath)
		require.NoError(t, err)
		defer f.Close()

		content, err := io.ReadAll(f)
		require.NoError(t, err)

		src := sourcecode.NewSourceCode(content, testCases[i].fileName)

		list, err := sqlChecker.Check(src)
		require.NoError(t, err)

		require.Len(t, list, len(testCases[i].lines))
		for j := range list {
			require.Equal(t, testCases[i].lines[j], list[j].Line)
		}
	}
}
