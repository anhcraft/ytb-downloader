package handle

import (
	"regexp"
	"strconv"
	"strings"
)

func extractPercentage(input string) (float64, bool) {
	// TODO optimize this
	if strings.HasPrefix(input, "[download]  ") &&
		strings.Contains(input, "of") &&
		strings.Contains(input, "at") {
		re := regexp.MustCompile(`\d*\.?\d*%`)
		match := re.FindStringSubmatch(input)
		if len(match) > 0 {
			if n, err := strconv.ParseFloat(match[0][0:len(match[0])-1], 64); err == nil {
				return n, true
			}
		}
	}
	return 0, false
}
