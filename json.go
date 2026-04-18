// Package letsgo provides generic useful functions for Go.
package letsgo

import "github.com/bytedance/sonic"

// DumpToJSON marshals v to a JSON string using bytedance/sonic, discarding
// any error. Intended for best-effort logging / debugging output where an
// unmarshalable value should not abort the calling code.
//
// Semantics:
//   - Pure with respect to v: no mutation of the input.
//   - Error-swallowing: any marshal error (including unsupported types like
//     chan, func, or complex numbers, or cycles) yields "" instead of an
//     error. Do NOT use this for on-the-wire or persistence payloads — use
//     sonic.Marshal / encoding/json directly so errors are surfaced.
//   - Returns "" when v is nil (avoiding the literal string "null" you would
//     otherwise get from the marshaller).
//   - Uses sonic.MarshalString, which avoids the []byte -> string copy of
//     sonic.Marshal, so this is cheaper than json.Marshal + string cast.
//   - Respects standard `json:"..."` struct tags.
//
// Example:
//
//	DumpToJSON(nil)                       // -> ""
//	DumpToJSON("hello")                   // -> `"hello"`
//	DumpToJSON(42)                        // -> "42"
//	DumpToJSON(map[string]int{"a": 1})    // -> `{"a":1}`
//	DumpToJSON(make(chan int))            // -> "" (unsupported type)
func DumpToJSON(v any) string {
	if v == nil {
		return ""
	}

	JSON, err := sonic.MarshalString(v)
	if err != nil {
		return ""
	}

	return JSON
}
