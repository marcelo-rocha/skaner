package runner_test

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/marcelo-rocha/skaner/domain/checker"
	"github.com/marcelo-rocha/skaner/domain/checker/exposure"
	"github.com/marcelo-rocha/skaner/domain/checker/sql"
	"github.com/marcelo-rocha/skaner/domain/checker/xss"
	"github.com/marcelo-rocha/skaner/domain/runner"
	"go.uber.org/zap"
)

type VulnerabilitiesSorter struct {
	list []checker.Vulnerability
}

func (s *VulnerabilitiesSorter) Len() int {
	return len(s.list)
}

func (s *VulnerabilitiesSorter) Swap(i, j int) {
	s.list[i], s.list[j] = s.list[j], s.list[i]
}

func (s *VulnerabilitiesSorter) Less(i, j int) bool {
	return strings.Compare(s.list[i].FilePath, s.list[j].FilePath) < 0 &&
		s.list[i].Line < s.list[j].Line &&
		strings.Compare(s.list[i].Kind, s.list[j].Kind) < 0
}

func TestRunExposure(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	options := runner.Options{
		SensitiveText: []string{"Switzerland National Bank", "Bill", "$2.2b"},
		JsonOutput:    true,
		WorkersQty:    2,
	}

	workDir, _ := os.Getwd()
	fileNames := []string{
		"exposure1.txt",
		"exposure2.txt",
		"exposure3.txt",
		"sql2.go",
		"xss1.html",
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
	require.Len(t, result, 6)

	sorter := VulnerabilitiesSorter{list: result}
	sort.Sort(&sorter)
	result = sorter.list

	require.Equal(t, 4, result[0].Line)

	require.Equal(t, exposure.VulnerabilityKind, result[0].Kind)
	require.Equal(t, exposure.VulnerabilityKind, result[1].Kind)
	require.Equal(t, exposure.VulnerabilityKind, result[2].Kind)
	require.Equal(t, sql.VulnerabilityKind, result[3].Kind)
	require.Equal(t, sql.VulnerabilityKind, result[4].Kind)
	require.Equal(t, xss.VulnerabilityKind, result[5].Kind)

	lines := checker.GetVulnerabilityLines(result)
	require.Equal(t, []int{4, 4, 5, 7, 15, 13}, lines)
}
