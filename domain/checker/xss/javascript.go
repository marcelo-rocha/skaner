package xss

import (
	"bytes"
	"io"

	"github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/js"
)

// ScanJS looks for a function call fnName in src, returning a slice of line numbers where
// the calls were found.
// This function returns all lines found even if a lexer error was found
func ScanJS(src []byte, fnName []byte) ([]int, error) {
	var result []int
	srcReader := bytes.NewReader(src)
	input := parse.NewInput(srcReader)
	lexer := js.NewLexer(input)

_scannerLoop:
	for {
		tt, text := lexer.Next()
		switch tt {
		case js.ErrorToken:
			if lexer.Err() != io.EOF {
				return result, NewParsingError("lexer error", lexer.Err())
			}
			break _scannerLoop
		case js.IdentifierToken:
			if bytes.Equal(text, fnName) {
				fnPos := input.Offset()
				if tt, _ := lexer.Next(); tt == js.OpenParenToken {
					if err := scanArguments(lexer); err != nil {
						return result, err
					}
					line, _, _ := parse.Position(bytes.NewReader(src), fnPos)
					result = append(result, line)
				}
			}
		}
	}
	return result, nil
}

func scanArguments(lexer *js.Lexer) error {
	parenthesesCount := 1
_scannerLoop:
	for {
		tt, _ := lexer.Next()
		switch tt {
		case js.ErrorToken:
			if lexer.Err() != io.EOF {
				return NewParsingError("lexer error", lexer.Err())
			}
			return ErrUnexpectedSourceEnd
		case js.OpenParenToken:
			parenthesesCount++
		case js.CloseParenToken:
			parenthesesCount--
			if parenthesesCount == 0 {
				break _scannerLoop
			}
		}
	}
	return nil
}
