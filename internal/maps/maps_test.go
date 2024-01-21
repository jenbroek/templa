package maps_test

import (
	"testing"

	. "github.com/jensbrks/templa/internal/maps"
	. "github.com/jensbrks/templa/internal/testing"

	"github.com/stretchr/testify/assert"
)

func ptr[T any](t T) *T {
	return &t
}

func TestMerge(t *testing.T) {
	type testCase struct {
		dst   map[string]any
		src   map[string]any
		want  map[string]any
		panic bool
	}

	RunTestCases(
		t,
		map[string]*testCase{
			"Adds new key": &testCase{
				dst:   map[string]any{"foo": "1"},
				src:   map[string]any{"bar": "2"},
				want:  map[string]any{"foo": "1", "bar": "2"},
				panic: false,
			},
			"Overwrites value of same type": &testCase{
				dst:   map[string]any{"foo": "1"},
				src:   map[string]any{"foo": "2"},
				want:  map[string]any{"foo": "2"},
				panic: false,
			},
			"Overwrites value of same pointer type": &testCase{
				dst:   map[string]any{"foo": ptr(1)},
				src:   map[string]any{"foo": ptr(2)},
				want:  map[string]any{"foo": ptr(2)},
				panic: false,
			},
			"Overwrites value of same array type": &testCase{
				dst:   map[string]any{"foo": [1]int{1}},
				src:   map[string]any{"foo": [1]int{2}},
				want:  map[string]any{"foo": [1]int{2}},
				panic: false,
			},
			"Does nothing with nil dest map": &testCase{
				dst:   nil,
				src:   map[string]any{"foo": "1"},
				want:  nil,
				panic: false,
			},
			"Does nothing with nil source map": &testCase{
				dst:   map[string]any{"foo": "1"},
				src:   nil,
				want:  map[string]any{"foo": "1"},
				panic: false,
			},
			"Does nothing with empty source map": &testCase{
				dst:   map[string]any{"foo": "1"},
				src:   map[string]any{},
				want:  map[string]any{"foo": "1"},
				panic: false,
			},
			"Merges array into assignable slice": &testCase{
				dst:   map[string]any{"nums": []any{"1"}},
				src:   map[string]any{"nums": [1]int{2}},
				want:  map[string]any{"nums": []any{"1", 2}},
				panic: false,
			},
			"Panics with non-assignable array types": &testCase{
				dst:   map[string]any{"nums": [1]string{"1"}},
				src:   map[string]any{"nums": [1]int{1}},
				want:  nil,
				panic: true,
			},
			"Merges assignable slice key": &testCase{
				dst:   map[string]any{"nums": []any{"1"}},
				src:   map[string]any{"nums": []int{2}},
				want:  map[string]any{"nums": []any{"1", 2}},
				panic: false,
			},
			"Merges assignable pointer slice key": &testCase{
				dst:   map[string]any{"nums": &[]any{"1"}},
				src:   map[string]any{"nums": &[]int{2}},
				want:  map[string]any{"nums": &[]any{"1", 2}},
				panic: false,
			},
			"Panics with non-assignable types": &testCase{
				dst:   map[string]any{"nums": "1"},
				src:   map[string]any{"nums": 1},
				want:  nil,
				panic: true,
			},
			"Panics with non-assignable slice types": &testCase{
				dst:   map[string]any{"nums": []string{"1"}},
				src:   map[string]any{"nums": []int{2}},
				want:  nil,
				panic: true,
			},
			"Merges assignable map key": &testCase{
				dst:   map[string]any{"nums": map[string]any{"1": "one"}},
				src:   map[string]any{"nums": map[string]int{"2": 2}},
				want:  map[string]any{"nums": map[string]any{"1": "one", "2": 2}},
				panic: false,
			},
			"Merges map key deeply": &testCase{
				dst:   map[string]any{"nums": map[string]any{"1": map[string]any{"en": "one"}}},
				src:   map[string]any{"nums": map[string]any{"1": map[string]any{"nl": "één"}}},
				want:  map[string]any{"nums": map[string]any{"1": map[string]any{"en": "one", "nl": "één"}}},
				panic: false,
			},
			"Panics with non-assignable map types": &testCase{
				dst:   map[string]any{"nums": map[string]string{"1": "one"}},
				src:   map[string]any{"nums": map[string]int{"2": 2}},
				want:  nil,
				panic: true,
			},
		},
		func(t *testing.T, tc *testCase) {
			if tc.panic {
				assert.Panics(t, func() { Merge(tc.dst, tc.src) })
			} else {
				if assert.NotPanics(t, func() { Merge(tc.dst, tc.src) }) {
					assert.Equal(t, tc.want, tc.dst)
				}
			}
		},
	)
}
