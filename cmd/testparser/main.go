package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/html"
)

// Tokenize HTML from stdin.
func main() {
	buf, _ := ioutil.ReadAll(os.Stdin)
	reader := bytes.NewReader(buf)
	input := parse.NewInput(reader)
	l := html.NewLexer(input)

	for {
		tt, data := l.Next()
		switch tt {
		case html.ErrorToken:
			if l.Err() != io.EOF {
				fmt.Println("Error on line :", l.Err())
			}
			return
		case html.StartTagToken:
			fmt.Println("Tag start", string(data))
			for {
				ttAttr, dataAttr := l.Next()
				if ttAttr != html.AttributeToken {
					break
				}

				key := dataAttr
				val := l.AttrVal()
				fmt.Println("Attribute", string(key), "=", string(val))
			}
			// ...
		case html.EndTagToken:
			fmt.Println("Tag end", string(data))
		case html.TextToken:
			line, _, con := parse.Position(bytes.NewReader(buf), input.Offset())
			fmt.Println("Tag text", string(data), " at", line, ":", con)

		}
	}
}
