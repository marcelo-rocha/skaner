package exposure

import (
	"bufio"
	"bytes"
	"errors"
	"strings"

	"github.com/marcelo-rocha/skaner/domain/checker"
)

const VulnerabilityKind = "Sensitive data exposure"

// TrimCuteset defines chars removed from the leading and trailing of a word
const TrimCuteset = ".,;:!?"

// ExposureChecker looks for strings in a unique line, case insensitive
// It can be used on multiple goroutines
type ExposureChecker struct {
	// list of strings case insensitive
	SensitiveData []string
	matchers      []StringMatcher
}

type CheckContext interface{}

func New(sensitiveData []string) (*ExposureChecker, error) {
	if len(sensitiveData) == 0 {
		return nil, errors.New("no sensitive data defined")
	}

	var matchers []StringMatcher
	for _, s := range sensitiveData {
		if strings.Contains(s, " ") {
			// It's a string sequence
			matchers = append(matchers, NewStringSequenceMatcher(s))
		} else {
			matchers = append(matchers, NewSimpleStringMatcher(s))
		}
	}

	return &ExposureChecker{
		SensitiveData: sensitiveData,
		matchers:      matchers,
	}, nil
}

func (c *ExposureChecker) Check(src checker.SourceCode) ([]checker.Vulnerability, error) {
	var result []checker.Vulnerability
	srcReader := bytes.NewReader(src.Bytes())
	srcScanner := bufio.NewScanner(srcReader)

	lineReader := bytes.NewReader([]byte{})

	currentLine := 1
	matcherCtxs := make([]MatcherContext, len(c.matchers))
	for i, m := range c.matchers {
		matcherCtxs[i] = m.BeginFile()
	}
	counters := make([]int, len(c.matchers))
	for srcScanner.Scan() {
		lineReader.Reset(srcScanner.Bytes())
		lineScanner := bufio.NewScanner(lineReader)
		lineScanner.Split(bufio.ScanWords)

		for i := range counters {
			counters[i] = 0
		}
		for lineScanner.Scan() {
			str := bytes.Trim(lineScanner.Bytes(), TrimCuteset)
			if len(str) > 0 {
				w := ToWord(str)
				for i, m := range c.matchers {
					if m.MatchWord(matcherCtxs[i], w) {
						counters[i] += 1
					}
				}
			}
		}
		if lineScanner.Err() != nil {
			return nil, lineScanner.Err()
		}

		found := true
		for i := range counters {
			found = found && counters[i] > 0
		}
		if found {
			result = append(result, checker.Vulnerability{
				Kind:     VulnerabilityKind,
				FilePath: src.FilePath(),
				Line:     currentLine,
			})
		}
		currentLine++
		for i, m := range c.matchers {
			m.NewLine(matcherCtxs[i])
		}
	}
	if srcScanner.Err() != nil {
		return nil, srcScanner.Err()
	}

	return result, nil
}

func (c *ExposureChecker) SupportedFileExtension(fileName string) bool {
	return true // all files
}
