package runner_test

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/marcelo-rocha/skaner/domain/checker"
	"github.com/marcelo-rocha/skaner/domain/runner"
	"go.uber.org/zap"
)

func TestRunExposure(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	options := runner.Options{
		SensitiveText:   []string{"Switzerland National Bank", "Bill", "$2.2b"},
		DisableXSSCheck: true,
		DisableSQLCheck: true,
		JsonOutput:      true,
		WorkersQty:      2,
	}

	workDir, _ := os.Getwd()
	fileNames := []string{
		"exposure1.txt",
		"exposure2.txt",
		"exposure3.txt",
	}
	var filePaths []string
	for _, n := range fileNames {
		filePaths = append(filePaths, path.Join(workDir, "../../test/data/", n))
	}

	buf := new(bytes.Buffer)

	runner.Run(context.Background(), buf, filePaths, options, logger)

	require.Greater(t, buf.Len(), 0)

	var result []checker.Vulnerability

	err := json.Unmarshal(buf.Bytes(), &result)
	require.NoError(t, err)
	require.Len(t, result, 2)
	//require.Equal(t, 4, result[0])
}
