package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

var (
	ctxFile     string
	globPattern stringSlice

	programName = filepath.Base(os.Args[0])

	helpFlag  = flag.Bool("h", false, "print this message and exit")
	listFlag  = flag.Bool("l", false, "list available context files")
	askFlag   = flag.Bool("i", false, "prompt before overwriting existing files")
	delimsStr = flag.String("d", "{{ }}", "left and right delimiter, separated by a space")
	ctxDir    = flag.String("c", getXdgDir(), "look in `directory` for context files")
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [-hl] | [-i] [-d delims] "+
		"[-c context_dir] [-p pattern ...] context_file\n", programName)
}

func run() int {
	bytes, err := ioutil.ReadFile(ctxFile)
	if err != nil {
		warn(err)
		return 1
	}

	ctx := map[string]interface{}{}
	if err := json.Unmarshal(bytes, &ctx); err != nil {
		warn(err)
		return 1
	}

	if len(globPattern) == 0 {
		file, err := os.Open(".templarc")
		if err != nil {
			warn(err)
			return 1
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			globPattern = append(globPattern, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			warn(err)
			return 1
		}
	}

	var paths []string
	for _, g := range globPattern {
		matches, err := filepath.Glob(g)
		if err != nil {
			warn(err)
			continue
		}

		for _, m := range matches {
			paths = append(paths, m)
		}
	}

	if len(paths) == 0 {
		return 0
	}

	delims := strings.Split(strings.Trim(*delimsStr, " "), " ")
	if len(delims) < 2 {
		warn("not enough delimiters specified")
		return 1
	}

	for _, path := range paths {
		fileInfo, err := os.Stat(path)
		if err != nil {
			warn(err)
			continue
		}
		if fileInfo.IsDir() {
			continue
		}

		// BUG
		// parsing multiple files in different directories
		// overwrites previous file with same basename
		tplName := filepath.Base(path)

		tpl, err := template.New(tplName).Delims(delims[0], delims[1]).ParseFiles(path)
		if err != nil {
			warn(err)
			continue
		}

		path := filepath.Join("..",
			strings.Join(
				strings.Split(path, string(os.PathSeparator))[1:],
				string(os.PathSeparator),
			),
		)

		if err := os.MkdirAll(filepath.Dir(path), 0777); err != nil {
			warn(err)
			continue
		}

		// BUG race condition between os.Stat and os.OpenFile
		if i, err := os.Stat(path); err == nil && i.IsDir() {
			warn("path '%s' is an existing directory, skipping")
			continue
		} else if *askFlag {
			fmt.Printf("file '%s' already exists, overwrite file? [n/Y] ", path)

			var ok string
			if _, err := fmt.Scanln(&ok); err != nil && err.Error() != "unexpected newline" {
				warn(err)
				continue
			}

			if !regexp.MustCompile("^[Yy].*$|^$").MatchString(ok) {
				continue
			}
		}

		file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fileInfo.Mode())
		if err != nil {
			warn(err)
			continue
		}
		defer file.Close()

		if err := tpl.ExecuteTemplate(file, tplName, ctx); err != nil {
			warn(err)
			continue
		}
	}

	return 0
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [-hl] | [-i] [-d delims] "+
		"[-c context_dir] [-p pattern ...] context_file\n", programName)
}

func init() {
	flag.Var(&globPattern, "p", "glob `pattern` of files to parse")

	flag.Usage = func() {
		usage()
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	if *ctxDir == getXdgDir() {
		if err := os.MkdirAll(*ctxDir, 0777); err != nil {
			warn(err)
			os.Exit(1)
		}
	}

	switch {
	case *helpFlag:
		flag.Usage()
		os.Exit(1)
	case *listFlag:
		files, _ := filepath.Glob(filepath.Join(*ctxDir, "*.json"))
		for _, f := range files {
			fmt.Println(strings.TrimSuffix(filepath.Base(f), ".json"))
		}
		os.Exit(0)
	case flag.NArg() == 0:
		usage()
		os.Exit(1)
	}

	if !strings.HasSuffix(flag.Arg(0), ".json") {
		ctxFile = filepath.Join(*ctxDir, flag.Arg(0)+".json")
	} else {
		ctxFile = filepath.Join(*ctxDir, flag.Arg(0))
	}

	os.Exit(run())
}
