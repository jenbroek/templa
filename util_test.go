package main_test

import (
	"reflect"
	"testing"

	"github.com/jensbrks/templa"
)

func TestMergeMaps(t *testing.T) {
	testCases := []struct {
		name string
		dst  map[string]any
		src  map[string]any
		want map[string]any
		fail bool
	}{
		{
			name: "Adds new key",
			dst:  map[string]any{"foo": "1"},
			src:  map[string]any{"bar": "2"},
			want: map[string]any{"foo": "1", "bar": "2"},
			fail: false,
		},
		{
			name: "Overwrites value of same non-map/slice type",
			dst:  map[string]any{"foo": "1"},
			src:  map[string]any{"foo": "2"},
			want: map[string]any{"foo": "2"},
			fail: false,
		},
		{
			name: "Overwrites value of different non-map/slice type",
			dst:  map[string]any{"foo": "1"},
			src:  map[string]any{"foo": 1},
			want: map[string]any{"foo": 1},
			fail: false,
		},
		{
			name: "Does nothing with nil dest map",
			dst:  nil,
			src:  map[string]any{"foo": "1"},
			want: nil,
			fail: false,
		},
		{
			name: "Does nothing with nil source map",
			dst:  map[string]any{"foo": "1"},
			src:  nil,
			want: map[string]any{"foo": "1"},
			fail: false,
		},
		{
			name: "Does nothing with empty source map",
			dst:  map[string]any{"foo": "1"},
			src:  map[string]any{},
			want: map[string]any{"foo": "1"},
			fail: false,
		},
		{
			name: "Merges slice key",
			dst:  map[string]any{"nums": []any{"1"}},
			src:  map[string]any{"nums": []any{"2"}},
			want: map[string]any{"nums": []any{"1", "2"}},
			fail: false,
		},
		{
			name: "Fails with non-slice and slice type",
			dst:  map[string]any{"nums": []any{"1"}},
			src:  map[string]any{"nums": "2"},
			want: nil,
			fail: true,
		},
		{
			name: "Fails with different slice types",
			dst:  map[string]any{"nums": []any{"1"}},
			src:  map[string]any{"nums": []int{2}},
			want: nil,
			fail: true,
		},
		{
			name: "Merges map key",
			dst:  map[string]any{"nums": map[string]any{"1": "one"}},
			src:  map[string]any{"nums": map[string]any{"2": "two"}},
			want: map[string]any{"nums": map[string]any{"1": "one", "2": "two"}},
			fail: false,
		},
		{
			name: "Merges map key deeply",
			dst:  map[string]any{"nums": map[string]any{"1": map[string]any{"en": "one"}}},
			src:  map[string]any{"nums": map[string]any{"1": map[string]any{"nl": "één"}}},
			want: map[string]any{"nums": map[string]any{"1": map[string]any{"en": "one", "nl": "één"}}},
			fail: false,
		},
		{
			name: "Fails with non-map and map type",
			dst:  map[string]any{"nums": map[string]any{"1": "one"}},
			src:  map[string]any{"nums": 2},
			want: nil,
			fail: true,
		},
		{
			name: "Fails with different map types",
			dst:  map[string]any{"nums": map[string]any{"1": "one"}},
			src:  map[string]any{"nums": map[string]int{"2": 2}},
			want: nil,
			fail: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if err := main.MergeMaps(tc.dst, tc.src); err != nil {
				if !tc.fail {
					t.Errorf("unexpected err: %v", err)
				}
			} else if !reflect.DeepEqual(tc.dst, tc.want) {
				t.Errorf("got: %v, want: %v", tc.dst, tc.want)
			}
		})
	}
}
