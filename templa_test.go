package main_test

import (
	"bytes"
	"cmp"
	"fmt"
	"reflect"
	"slices"
	"testing"
	"testing/fstest"
	"text/template"

	. "github.com/jensbrks/templa"
)

func TestRun(t *testing.T) {
	var wr bytes.Buffer
	fsys := fstest.MapFS{
		"greeting":  {Data: []byte("hello {{ .name }}")},
		"data.yaml": {Data: []byte("name: Bob")},
	}
	tmplPaths := []string{"greeting"}
	valueFiles := []string{"data.yaml"}

	want := "hello Bob"

	if err := Run(&wr, fsys, tmplPaths, valueFiles); err != nil {
		ErrUnexpected(t, err)
	} else if got := wr.String(); got != want {
		ErrNotEqual(t, got, want)
	}
}

func TestParseTemplates(t *testing.T) {
	fsys := fstest.MapFS{
		"hello": {Data: []byte("hello")},
		"bye":   {Data: []byte("bye")},
	}
	tmplPaths := []string{"hello", "bye"}

	got, err := ParseTemplates(fsys, tmplPaths)

	if err != nil {
		ErrUnexpected(t, err)
	} else {
		if got == nil {
			ErrNotEqual(t, got, "not nil")
		} else {
			// Sort to avoid flaky tests
			tmpls := got.Templates()
			slices.SortFunc(tmpls, func(a, b *template.Template) int {
				return cmp.Compare(a.Name(), b.Name())
			})
			slices.Sort(tmplPaths)

			for i, tp := range tmpls {
				want := tmplPaths[i]
				if tp.Name() != want {
					ErrNotEqual(
						t,
						fmt.Sprintf("name: %d: %s", i, tp.Name()),
						fmt.Sprintf("name: %d: %s", i, want),
					)
				}
			}
		}
	}
}

func TestReadValueFiles(t *testing.T) {
	type testCase struct {
		fsys       fstest.MapFS
		valueFiles []string
		want       map[string]any
	}

	RunTestCases(
		t,
		map[string]*testCase{
			"Parses values from both YAML inputs": &testCase{
				fsys: fstest.MapFS{
					"foo": {Data: []byte("{foo: bar}")},
					"bar": {Data: []byte("{bar: baz}")},
				},
				valueFiles: []string{"foo", "bar"},
				want:       map[string]any{"foo": "bar", "bar": "baz"},
			},
			"Merges lists from both YAML inputs": &testCase{
				fsys: fstest.MapFS{
					"nums12": {Data: []byte("{nums: [1,2]}")},
					"nums3":  {Data: []byte("{nums: [3]}")},
				},
				valueFiles: []string{"nums12", "nums3"},
				want:       map[string]any{"nums": []any{1, 2, 3}},
			},
			"Merges maps from both YAML inputs": &testCase{
				fsys: fstest.MapFS{
					"nums_one": {Data: []byte("{nums: {'1': one}}")},
					"nums_two": {Data: []byte("{nums: {'2': two}}")},
				},
				valueFiles: []string{"nums_one", "nums_two"},
				want:       map[string]any{"nums": map[string]any{"1": "one", "2": "two"}},
			},
		},
		func(t *testing.T, tc *testCase) {
			got, err := ReadValueFiles(tc.fsys, tc.valueFiles)
			if err != nil {
				ErrUnexpected(t, err)
			} else if !reflect.DeepEqual(got, tc.want) {
				ErrNotEqual(t, got, tc.want)
			}
		},
	)
}
