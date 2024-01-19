package main

import (
	"bytes"
	"testing"
	"testing/fstest"
	"text/template"

	. "github.com/jensbrks/templa/internal/testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
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
	err := run(&wr, fsys, tmplPaths, valueFiles)
	got := wr.String()

	if assert.NoError(t, err) {
		assert.Equal(t, want, got)
	}
}

func TestParseTemplates(t *testing.T) {
	fsys := fstest.MapFS{
		"hello": {Data: []byte("hello")},
		"bye":   {Data: []byte("bye")},
	}
	tmplPaths := []string{"hello", "bye"}

	got, err := parseTemplates(fsys, tmplPaths)

	if assert.NoError(t, err) {
		tmplNames := lo.Map(got.Templates(), func(t *template.Template, _ int) string {
			return t.Name()
		})

		assert.ElementsMatch(t, tmplPaths, tmplNames)
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
			got, err := readValueFiles(tc.fsys, tc.valueFiles)
			if assert.NoError(t, err) {
				assert.Equal(t, tc.want, got)
			}
		},
	)
}
