package checker

import "errors"

type SourceCode interface {
	Bytes() []byte
	FilePath() string
}

type Vulnerability struct {
	Kind     string `json:"kind"`
	FilePath string `json:"filePath"`
	Line     int    `json:"lineNumber"`
}

type Checker interface {
	Check(src SourceCode) ([]Vulnerability, error)
	SupportedFileExtension(fileName string) bool
}

type NoOperationChecker struct{}

func (c *NoOperationChecker) Check(src SourceCode) ([]Vulnerability, error) {
	return nil, errors.New("no security checking defined")
}

func (c *NoOperationChecker) SupportedFileExtension(fileName string) bool {
	return false
}

func GetVulnerabilityLines(list []Vulnerability) []int {
	var r []int
	for i := range list {
		r = append(r, list[i].Line)
	}
	return r
}
