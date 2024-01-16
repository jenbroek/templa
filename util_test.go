package main_test

import (
	"reflect"
	"testing"

	. "github.com/jensbrks/templa"
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
			"Overwrites value of same non-map/slice type": &testCase{
				dst:  map[string]any{"foo": "1"},
				src:  map[string]any{"foo": "2"},
				want: map[string]any{"foo": "2"},
				fail: false,
			},
			"Overwrites value of different non-map/slice type": &testCase{
				dst:  map[string]any{"foo": "1"},
				src:  map[string]any{"foo": 1},
				want: map[string]any{"foo": 1},
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
			"Merges slice key": &testCase{
				dst:  map[string]any{"nums": []any{"1"}},
				src:  map[string]any{"nums": []any{"2"}},
				want: map[string]any{"nums": []any{"1", "2"}},
				fail: false,
			},
			"Fails with non-slice and slice type": &testCase{
				dst:  map[string]any{"nums": []any{"1"}},
				src:  map[string]any{"nums": "2"},
				want: nil,
				fail: true,
			},
			"Fails with different slice types": &testCase{
				dst:  map[string]any{"nums": []any{"1"}},
				src:  map[string]any{"nums": []int{2}},
				want: nil,
				fail: true,
			},
			"Merges map key": &testCase{
				dst:  map[string]any{"nums": map[string]any{"1": "one"}},
				src:  map[string]any{"nums": map[string]any{"2": "two"}},
				want: map[string]any{"nums": map[string]any{"1": "one", "2": "two"}},
				fail: false,
			},
			"Merges map key deeply": &testCase{
				dst:  map[string]any{"nums": map[string]any{"1": map[string]any{"en": "one"}}},
				src:  map[string]any{"nums": map[string]any{"1": map[string]any{"nl": "één"}}},
				want: map[string]any{"nums": map[string]any{"1": map[string]any{"en": "one", "nl": "één"}}},
				fail: false,
			},
			"Fails with non-map and map type": &testCase{
				dst:  map[string]any{"nums": map[string]any{"1": "one"}},
				src:  map[string]any{"nums": 2},
				want: nil,
				fail: true,
			},
			"Fails with different map types": &testCase{
				dst:  map[string]any{"nums": map[string]any{"1": "one"}},
				src:  map[string]any{"nums": map[string]int{"2": 2}},
				want: nil,
				fail: true,
			},
		},
		func(t *testing.T, tc *testCase) {
			if err := MergeMaps(tc.dst, tc.src); err != nil {
				if !tc.fail {
					ErrUnexpected(t, err)
				}
			} else if !reflect.DeepEqual(tc.dst, tc.want) {
				ErrNotEqual(t, tc.dst, tc.want)
			}
		},
	)
}
