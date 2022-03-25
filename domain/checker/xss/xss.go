package xss

import (
	"bytes"
	"io"
	"path"

	"github.com/marcelo-rocha/skaner/domain/checker"
	"github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/html"
)

const VulnerabilityKind = "Cross site scripting"

var alertText = []byte("alert")

type XSSChecker struct {
}

func New() *XSSChecker {
	return &XSSChecker{}
}

func (c *XSSChecker) Check(src checker.SourceCode) ([]checker.Vulnerability, error) {
	ext := path.Ext(src.FilePath())
	switch ext {
	case ".js":
		return checkJS(src)

	case ".htm", ".html":
		return checkHTML(src)
	}
	return nil, ErrFileTypeNotSupported
}

func (c *XSSChecker) SupportedFileExtension(fileName string) bool {
	ext := path.Ext(fileName)
	return ext == ".js" || ext == ".htm" || ext == ".html"
}

var (
	beginScriptTag = []byte("<script")
	endScriptTag   = []byte("</script>")
)

func checkHTML(src checker.SourceCode) ([]checker.Vulnerability, error) {
	var result []checker.Vulnerability
	reader := bytes.NewReader(src.Bytes())
	input := parse.NewInput(reader)
	lexer := html.NewLexer(input)
	scriptTagLevel := 0
	scriptTagOffset := 0

_parserLoop:
	for {
		tt, data := lexer.Next()
		switch tt {
		case html.ErrorToken:
			if lexer.Err() != io.EOF {
				return result, NewParsingError("html lexer error", lexer.Err())
			}
			break _parserLoop
		case html.StartTagToken:
			if bytes.EqualFold(data, beginScriptTag) {
				scriptTagLevel++
				scriptTagOffset = input.Offset()
			}
		case html.EndTagToken:
			if bytes.EqualFold(data, endScriptTag) {
				scriptTagLevel--
				if scriptTagLevel < 0 {
					return result, ErrUnbalancedScriptTag
				}
			}
		case html.TextToken:
			if scriptTagLevel > 0 {
				// ScanJS errors are ignored aiming to cover all the file
				lines, _ := ScanJS(data, alertText)
				if len(lines) > 0 {
					baseLine, _, _ := parse.Position(bytes.NewReader(src.Bytes()), scriptTagOffset)
					for _, line := range lines {
						result = append(result, checker.Vulnerability{
							Kind:     VulnerabilityKind,
							FilePath: src.FilePath(),
							Line:     line + baseLine - 1,
						})
					}
				}
			}
		}
	}
	return result, nil
}

func checkJS(src checker.SourceCode) ([]checker.Vulnerability, error) {
	var result []checker.Vulnerability
	lines, err := ScanJS(src.Bytes(), alertText)
	if len(lines) > 0 {
		for _, line := range lines {
			result = append(result, checker.Vulnerability{
				Kind:     VulnerabilityKind,
				FilePath: src.FilePath(),
				Line:     line,
			})
		}
	}

	return result, err
}
