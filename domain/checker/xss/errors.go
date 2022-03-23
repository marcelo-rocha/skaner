package xss

import (
	"errors"
)

var (
	ErrInvalidSource        = errors.New("invalid source code")
	ErrFileTypeNotSupported = errors.New("file type not supported")
	ErrUnbalancedScriptTag  = errors.New("unbalanced script tag")
	ErrUnexpectedSourceEnd  = errors.New("unexpected source end")
)

type ParsingError struct {
	Msg string
	Err error
}

func NewParsingError(msg string, originalError error) error {
	return &ParsingError{
		Msg: msg,
		Err: originalError,
	}
}

func (p *ParsingError) Error() string {
	return p.Msg
}

func (p *ParsingError) Unwrap() error {
	return p.Err
}
