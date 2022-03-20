package xss

import (
	"github.com/marcelo-rocha/skaner/domain/checker"
	"github.com/marcelo-rocha/skaner/domain/sourcecode"
)

type XSSChecker struct {
}

func New() *XSSChecker {
	return &XSSChecker{}
}

func (c *XSSChecker) Check(src *sourcecode.SourceCode) ([]checker.Vulnerability, error) {
	return nil, nil
}
