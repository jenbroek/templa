package main

import (
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/jensbrks/templa/internal/maps"

	"github.com/spf13/pflag"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"
)

// See: https://github.com/golang/go/issues/44286
type osFS struct{}

func (*osFS) Open(name string) (fs.File, error) { return os.Open(name) }

var (
	VERSION     string
	programName = filepath.Base(os.Args[0])

	openDelim   = pflag.StringP("open-delim", "o", "{{", "Define the opening delimiter")
	closeDelim  = pflag.StringP("close-delim", "c", "}}", "Define the closing delimiter")
	valueFiles  = pflag.StringSliceP("values", "f", []string{}, "List of value files to use for templating")
	versionFlag = pflag.CountP("version", "v", "Print version information to stderr and exit")
)

func init() {
	log.SetPrefix(programName + ": ")
	// Clear flags
	log.SetFlags(0)
}

func main() {
	pflag.Parse()

	if *versionFlag > 0 {
		log.Printf("%s-%s\n", programName, VERSION)
		return
	}

	var tmplPaths []string

	if pflag.NArg() == 0 {
		tmplPaths = []string{os.Stdin.Name()}
	} else {
		tmplPaths = pflag.Args()
	}

	if err := run(os.Stdout, new(osFS), tmplPaths, *valueFiles); err != nil {
		log.Fatal(err)
	}
}

func run(wr io.Writer, fsys fs.FS, tmplPaths, valueFiles []string) error {
	var g errgroup.Group
	var tmpl *template.Template
	var data map[string]any

	g.Go(func() (err error) {
		tmpl, err = parseTemplates(fsys, tmplPaths)
		return
	})

	g.Go(func() (err error) {
		data, err = readValueFiles(fsys, valueFiles)
		return
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return tmpl.ExecuteTemplate(wr, tmpl.Name(), data)
}

func parseTemplates(fsys fs.FS, tmplPaths []string) (*template.Template, error) {
	var parentTmpl *template.Template

	for _, tmplPath := range tmplPaths {
		// `template#Template.ParseFiles` forces the template name to be the basename
		// of the specified path(s). In order to use the full (relative) path, we
		// must call `template#Template.Parse` ourselves.
		bytes, err := fs.ReadFile(fsys, tmplPath)
		if err != nil {
			return nil, err
		}

		var tmpl *template.Template
		if parentTmpl == nil {
			parentTmpl = template.New(tmplPath)
			tmpl = parentTmpl
		} else {
			tmpl = parentTmpl.New(tmplPath)
		}

		if _, err = tmpl.Delims(*openDelim, *closeDelim).Option("missingkey=zero").Parse(string(bytes)); err != nil {
			return nil, err
		}
	}

	return parentTmpl, nil
}

func readValueFiles(fsys fs.FS, valueFiles []string) (map[string]any, error) {
	data := make(map[string]any)

	for _, fp := range valueFiles {
		bytes, err := fs.ReadFile(fsys, fp)
		if err != nil {
			return nil, err
		}

		m := make(map[string]any)
		if err = yaml.Unmarshal(bytes, m); err != nil {
			return nil, err
		}

		maps.Merge(data, m)
	}

	return data, nil
}
