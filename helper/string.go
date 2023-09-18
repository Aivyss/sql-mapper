package helper

import "strings"

func ReplaceNewLineAndTabToSpace(input string) string {
	result := strings.ReplaceAll(input, "\n", " ")
	result = strings.ReplaceAll(result, "\r", " ")
	result = strings.ReplaceAll(result, "\t", " ")
	result = strings.TrimSpace(result)
	return result
}

func IsBlank(input string) bool {
	input = ReplaceNewLineAndTabToSpace(input)
	input = strings.TrimSpace(input)

	return len(input) == 0
}
