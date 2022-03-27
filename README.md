# skaner

** This code is a Code Challenger, it's not production ready **

Skaner is a security code scanner. 

It looks for three types of vulnerabilities: 
- Exposure of sensitive text
- Cross site scripting
- SQL Injection

## Instructions

It's necessary Go 1.17 or more recent

To build, execute `make build`. The executable is saved on `dist` folder.

To run tests, execute `make test`.

Execute `skaner` to see the instructions.

The exposure checker requires the list of exposures strings. If we don't look for exposures,
inform the parameter `--no-exposure-checker`

Below some examples of how to use it:

```bash
./dist/skaner --no-exposure-checker test/data/sql1.go 
./dist/skaner --no-exposure-checker test/data/xss1.html
./dist/skaner --sensitive-text="Bill,\$2.2b,Switzerland" test/data/exposure2.txt

```
Notice that it's necessary to escape the $ symbol.

We could pass multiple files:

```bash
./dist/skaner --sensitive-text="Bill,\$2.2b,Switzerland" test/data/exposure2.txt test/data/sql1.go

```

The default output format is plain text, but skaner support JSON output:

```bash
./dist/skaner --no-exposure-checker --json test/data/xss1.html
```


## Code design

Each vulnerability is verified by an object that implements the Checker interface:

```go
type SourceCode interface {
	Bytes() []byte
	FilePath() string
}

type Vulnerability struct {
	Kind     string `json:"kind"`
	FilePath string `json:"filePath"`
	Line     int    `json:"lineNumber"`
}

type Checker interface {
	Check(src SourceCode) ([]Vulnerability, error)
	SupportedFileExtension(fileName string) bool
}
```

All checkers objects support be called by multiple goroutines concurrently. They uses as most as possible slices, to prevent unecessary memory copy.

## Implemented Checkers

### Sensitive Text Exposure Checker
It looks for multiple strings at the same line. It read the source line by line and compare the words. It's used a fast hash function for comparison of words.

### SQL Injection Checker
It uses a regular expression for loop for quoted strings, then it scan the strings looking for a pattern 
`...SELECT...WHERE...%s`. If the pattern is found, another regular expression is used to find the vulnerabilities line numbers.

### Cross site scripting checker
It uses lexers for HTML and Javascript. When a _script_ tag is found in a HTML content, the script is scanned using
the Javascript lexer looking for an `alert()` function call.

## Runner
Runner package uses multiples goroutines to check for vulnerabilities concurrently




