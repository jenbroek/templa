package main

import "testing"

var (
	Run            = run
	ParseTemplates = parseTemplates
	ReadValueFiles = readValueFiles
)

func RunTestCases[T any](t *testing.T, tcs map[string]T, fn func(t *testing.T, tc T)) {
	for n, tc := range tcs {
		// capture range variables
		n, tc := n, tc
		t.Run(n, func(t *testing.T) {
			t.Parallel()
			fn(t, tc)
		})
	}
}
