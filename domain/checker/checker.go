package checker

import (
	"github.com/marcelo-rocha/skaner/domain/sourcecode"
)

type Vulnerability struct {
	Kind     string
	FilePath string
	Line     int
}

type Checker interface {
	Check(src *sourcecode.SourceCode) ([]Vulnerability, error)
}
