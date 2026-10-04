package conversion

import (
	"bytes"
	"regexp"
)

var tagChecker = regexp.MustCompile("<.*?>")
var whitespaceChecker = regexp.MustCompile(`\s{2,}`)

func StripTagsFromHtml(input []byte) []byte {
	output := tagChecker.ReplaceAll(input, []byte{})
	output = bytes.ReplaceAll(output, []byte("\n"), []byte(" "))
	output = bytes.ReplaceAll(output, []byte("	"), []byte(" "))
	output = whitespaceChecker.ReplaceAll(output, []byte(" "))
	return output
}
