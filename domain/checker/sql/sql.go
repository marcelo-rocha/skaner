package sql

import (
	"github.com/marcelo-rocha/skaner/domain/checker"
	"github.com/marcelo-rocha/skaner/domain/sourcecode"
)

type SQLChecker struct {
}

func New() *SQLChecker {
	return &SQLChecker{}
}

func (c *SQLChecker) Check(src *sourcecode.SourceCode) ([]checker.Vulnerability, error) {
	return nil, nil
}
