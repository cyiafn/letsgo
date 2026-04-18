// Package letsgo provides generic useful functions for Go.
package letsgo

import "encoding/json"

// DumpToJSON marshals v to a JSON string using encoding/json, discarding
// any error. Intended for best-effort logging / debugging output where an
// unmarshalable value should not abort the calling code.
//
// Semantics:
//   - Error-swallowing: any marshal error (including unsupported types like
//     chan, func, or complex numbers, or cycles) yields "" instead of an
//     error. Do NOT use this for on-the-wire or persistence payloads — call
//     json.Marshal directly so errors are surfaced.
//   - Returns "" when v is nil (avoiding the literal string "null" you would
//     otherwise get from the marshaller). Note: this is an interface-nil
//     check; a typed nil (e.g. a (*T)(nil) stored in an `any`) is NOT
//     caught here and will be marshalled as "null".
//   - Uses standard encoding/json behavior: HTML characters (<, >, &) are
//     escaped, and map keys are emitted in sorted order — so output is
//     deterministic across runs.
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

	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}

	return string(b)
}
