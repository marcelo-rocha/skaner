package sourcecode

type SourceCode struct {
	Source   []byte
	FilePath string
}

// NewSourceCode represents a program source code. src should be utf-8 format
func NewSourceCode(src []byte, path string) *SourceCode {
	return &SourceCode{
		Source:   src,
		FilePath: path,
	}
}
