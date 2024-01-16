package main

import "testing"

var MergeMaps = mergeMaps

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

func ErrNotEqual[T any](t *testing.T, got, want T) {
	t.Errorf("got: %v, want %v", got, want)
}
