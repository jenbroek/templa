package main

import "testing"

var (
	Run            = run
	ParseTemplates = parseTemplates
	ReadValueFiles = readValueFiles
	MergeMaps      = mergeMaps
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

func ErrUnexpected(t *testing.T, err error) {
	t.Errorf("unexpected err: %v", err)
}

func ErrNotEqual[T1, T2 any](t *testing.T, got T1, want T2) {
	t.Errorf("got: %v, want %v", got, want)
}
