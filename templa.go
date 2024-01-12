package main

import (
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

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

	var tmplPaths []string

	if *versionFlag > 0 {
		log.Printf("%s-%s\n", programName, VERSION)
		return
	}

	if pflag.NArg() == 0 {
		tmplPaths = []string{os.Stdin.Name()}
	} else {
		tmplPaths = pflag.Args()
	}

	if err := run(tmplPaths); err != nil {
		log.Fatal(err)
	}
}

func run(tmplPaths []string) error {
	tmpl, err := parseTemplates(tmplPaths)
	if err != nil {
		return err
	}

	data, err := readValueFiles()
	if err != nil {
		return err
	}

	return tmpl.ExecuteTemplate(os.Stdout, tmpl.Name(), data)
}

func parseTemplates(tmplPaths []string) (*template.Template, error) {
	var parentTmpl *template.Template

	for _, tmplPath := range tmplPaths {
		// `template#Template.ParseFiles` forces the template name to be the basename
		// of the specified path(s). In order to use the full (relative) path, we
		// must call `template#Template.Parse` ourselves.
		bytes, err := os.ReadFile(tmplPath)
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

func readValueFiles() (map[string]any, error) {
	data := make(map[string]any)

	for _, valueFile := range *valueFiles {
		bytes, err := os.ReadFile(valueFile)
		if err != nil {
			return nil, err
		}

		m := make(map[string]any)
		if err = yaml.Unmarshal(bytes, m); err != nil {
			return nil, err
		}

		if err = mergeMaps(data, m); err != nil {
			return nil, err
		}
	}

	return data, nil
}
