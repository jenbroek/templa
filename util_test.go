package main_test

import (
	"testing"

	. "github.com/jensbrks/templa"

	"github.com/stretchr/testify/assert"
)

func TestMergeMaps(t *testing.T) {
	type testCase struct {
		dst  map[string]any
		src  map[string]any
		want map[string]any
		fail bool
	}

	RunTestCases(
		t,
		map[string]*testCase{
			"Adds new key": &testCase{
				dst:  map[string]any{"foo": "1"},
				src:  map[string]any{"bar": "2"},
				want: map[string]any{"foo": "1", "bar": "2"},
				fail: false,
			},
			"Overwrites value of same type": &testCase{
				dst:  map[string]any{"foo": "1"},
				src:  map[string]any{"foo": "2"},
				want: map[string]any{"foo": "2"},
				fail: false,
			},
			"Does nothing with nil dest map": &testCase{
				dst:  nil,
				src:  map[string]any{"foo": "1"},
				want: nil,
				fail: false,
			},
			"Does nothing with nil source map": &testCase{
				dst:  map[string]any{"foo": "1"},
				src:  nil,
				want: map[string]any{"foo": "1"},
				fail: false,
			},
			"Does nothing with empty source map": &testCase{
				dst:  map[string]any{"foo": "1"},
				src:  map[string]any{},
				want: map[string]any{"foo": "1"},
				fail: false,
			},
			"Merges assignable slice key": &testCase{
				dst:  map[string]any{"nums": []any{"1"}},
				src:  map[string]any{"nums": []int{2}},
				want: map[string]any{"nums": []any{"1", 2}},
				fail: false,
			},
			"Fails with non-assignable types": &testCase{
				dst:  map[string]any{"nums": "1"},
				src:  map[string]any{"nums": 1},
				want: nil,
				fail: true,
			},
			"Fails with non-assignable slice types": &testCase{
				dst:  map[string]any{"nums": []string{"1"}},
				src:  map[string]any{"nums": []int{2}},
				want: nil,
				fail: true,
			},
			"Merges assignable map key": &testCase{
				dst:  map[string]any{"nums": map[string]any{"1": "one"}},
				src:  map[string]any{"nums": map[string]int{"2": 2}},
				want: map[string]any{"nums": map[string]any{"1": "one", "2": 2}},
				fail: false,
			},
			"Merges map key deeply": &testCase{
				dst:  map[string]any{"nums": map[string]any{"1": map[string]any{"en": "one"}}},
				src:  map[string]any{"nums": map[string]any{"1": map[string]any{"nl": "één"}}},
				want: map[string]any{"nums": map[string]any{"1": map[string]any{"en": "one", "nl": "één"}}},
				fail: false,
			},
			"Fails with non-assignable map types": &testCase{
				dst:  map[string]any{"nums": map[string]string{"1": "one"}},
				src:  map[string]any{"nums": map[string]int{"2": 2}},
				want: nil,
				fail: true,
			},
		},
		func(t *testing.T, tc *testCase) {
			err := MergeMaps(tc.dst, tc.src)
			if tc.fail {
				assert.Error(t, err)
			} else {
				if assert.NoError(t, err) {
					assert.Equal(t, tc.want, tc.dst)
				}
			}
		},
	)
}
