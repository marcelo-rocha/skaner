package checker

type SourceCode interface {
	Bytes() []byte
	FilePath() string
}

type Vulnerability struct {
	Kind     string
	FilePath string
	Line     int
}

type Checker interface {
	Check(src SourceCode) ([]Vulnerability, error)
	SupportedFileExtension(fileName string) bool
}
