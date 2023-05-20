// Package letsgo provides generic useful functions for Go.
package letsgo

import "encoding/json"

// DumpToJSON dumps the given value to JSON format and returns a string only.
// This is useful for logging purposes where you want to log the value in JSON format without an error.
// It fails silently and returns an empty string if it fails to dump the value to JSON.
func DumpToJSON(v any) string {
	if v == nil {
		return ""
	}

	JSON, err := json.Marshal(v)
	if err != nil {
		return ""
	}

	return string(JSON)
}
