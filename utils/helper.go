package utils

import (
	"fmt"
	"regexp"
	"strings"
)

var shellMetaPattern = regexp.MustCompile(`[&|><^]`)

func NormalizeText(input string) string {
	fields := strings.Fields(strings.TrimSpace(strings.ToLower(input)))
	return strings.Join(fields, " ")
}

func QuoteCMDArg(value string) string {
	escaped := strings.ReplaceAll(value, `"`, `""`)
	return fmt.Sprintf(`"%s"`, escaped)
}

func ContainsShellMetacharacters(value string) bool {
	return shellMetaPattern.MatchString(value)
}
