package lib

import (
	"regexp"
)

var g_isVarNameRegexp *regexp.Regexp = regexp.MustCompile(`^[a-zA-Z_]\w*$`)
var g_fetchListTypeRegexp *regexp.Regexp = regexp.MustCompile(`^list{(.+)}$`)
var g_camelToUnderscoreCase1Regexp *regexp.Regexp = regexp.MustCompile(`([A-Z][0-9]*)([A-Z][0-9]*[a-z])`)
var g_camelToUnderscoreCase2Regexp *regexp.Regexp = regexp.MustCompile(`([a-z][0-9]*)([A-Z])`)
var g_notWordRegexp *regexp.Regexp = regexp.MustCompile(`[^\w]`)
