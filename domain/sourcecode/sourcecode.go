package sourcecode

// SourceCode represents a program source code. src should be utf-8 format
// It implements checker.Source interface
type SourceCode struct {
	source   []byte
	filePath string
}

func NewSourceCode(src []byte, path string) *SourceCode {
	return &SourceCode{
		source:   src,
		filePath: path,
	}
}

func (s *SourceCode) Bytes() []byte {
	return s.source
}

func (s *SourceCode) FilePath() string {
	return s.filePath
}
