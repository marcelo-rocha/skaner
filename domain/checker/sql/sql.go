package sql

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"regexp"

	"github.com/marcelo-rocha/skaner/domain/checker"
)

const VulnerabilityKind = "SQL Injection"

type SQLChecker struct {
	quotedStrRE *regexp.Regexp
	lineBreakRE *regexp.Regexp
}

func New() *SQLChecker {
	return &SQLChecker{
		quotedStrRE: regexp.MustCompile(`"([^"\\\n]|\\\n|\\\w)*"`),
		lineBreakRE: regexp.MustCompile(`\n`),
	}
}

func (c *SQLChecker) Check(src checker.SourceCode) ([]checker.Vulnerability, error) {
	strIndexes := c.quotedStrRE.FindAllIndex(src.Bytes(), -1)
	if len(strIndexes) == 0 {
		return []checker.Vulnerability{}, nil
	}

	strReader := bytes.NewReader([]byte{})
	var insecureOffsets []int
	for _, p := range strIndexes {
		str := src.Bytes()[p[0]+1 : p[1]-1]
		strReader.Reset(str)
		if c.checkInsecureSQL(strReader) {
			insecureOffsets = append(insecureOffsets, p[0])
		}
	}

	if len(insecureOffsets) == 0 {
		return []checker.Vulnerability{}, nil
	}

	result := make([]checker.Vulnerability, len(insecureOffsets))
	lines := c.resolveLineNumbers(src, insecureOffsets)
	if len(lines) != len(result) {
		return nil, errors.New("unexpected error")
	}
	for i, n := range lines {
		result[i] = checker.Vulnerability{
			Kind:     VulnerabilityKind,
			FilePath: src.FilePath(),
			Line:     n,
		}
	}

	return result, nil
}

var (
	selectKeyword    = []byte("SELECT")
	whereKeyword     = []byte("WHERE")
	placeholdKeyword = []byte("%s")
)

type checkAutomataState int

const (
	StringBegin checkAutomataState = iota
	SelectFound
	WhereFound
	PlaceHolderFound
)

// checkInsecureSQL looks if a string is in the format ...SELECT...WHERE...%s
func (c *SQLChecker) checkInsecureSQL(strReader io.Reader) bool {
	srcScanner := bufio.NewScanner(strReader)
	srcScanner.Split(bufio.ScanWords)

	state := StringBegin
	for srcScanner.Scan() {
		w := srcScanner.Bytes()
		switch state {
		case StringBegin:
			if bytes.EqualFold(w, selectKeyword) {
				state = SelectFound
			}
		case SelectFound:
			if bytes.EqualFold(w, whereKeyword) {
				state = WhereFound
			}
		case WhereFound:
			if bytes.Equal(w, placeholdKeyword) {
				return true
			}
		}
	}

	return false
}

// resolveLineNumbers offsets argument must be ordered
func (c *SQLChecker) resolveLineNumbers(src checker.SourceCode, offsets []int) []int {
	result := make([]int, len(offsets))
	reader := bytes.NewReader(src.Bytes())
	line := 0

_offsets_loop:
	for i, n := range offsets {
		for {
			line++
			idx := c.lineBreakRE.FindReaderIndex(reader)
			if idx == nil {
				break _offsets_loop
			}
			if idx[0] > n {
				break
			}
		}
		result[i] = line
	}

	return result
}
