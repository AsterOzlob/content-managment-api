package utils

import "github.com/microcosm-cc/bluemonday"

var Policy = bluemonday.UGCPolicy()

func Sanitize(s string) string {
	return Policy.Sanitize(s)
}
