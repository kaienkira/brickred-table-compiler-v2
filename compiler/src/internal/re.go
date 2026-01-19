package lib

import (
	"regexp"
)

var g_isVarNameRegexp *regexp.Regexp = regexp.MustCompile(`^[a-zA-Z_]\w*$`)
