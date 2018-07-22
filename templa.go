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

type stringSlice []string

func (sl *stringSlice) String() string {
	return fmt.Sprint(*sl)
}

func (sl *stringSlice) Set(value string) error {
	*sl = append(*sl, value)
	return nil
}

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

func init() {
	flag.Var(&globPattern, "p", "glob `pattern` of files to parse")

	flag.Usage = func() {
		usage()
		flag.PrintDefaults()
		os.Exit(1)
	}
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

	var files []string
	for _, g := range globPattern {
		matches, err := filepath.Glob(g)
		if err != nil {
			warn(err)
			continue
		}

		for _, m := range matches {
			files = append(files, m)
		}
	}

	if len(files) == 0 {
		return 0
	}

	delims := strings.Split(strings.Trim(*delimsStr, " "), " ")
	if len(delims) < 2 {
		warn("not enough delimiters specified")
		return 1
	}

	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			warn(err)
			continue
		}
		if info.IsDir() {
			continue
		}

		// BUG
		// parsing multiple files in different directories
		// overwrites previous file with same basename
		tplName := filepath.Base(file)

		tpl, err := template.New(tplName).Delims(delims[0], delims[1]).ParseFiles(file)
		if err != nil {
			warn(err)
			continue
		}

		destPath := filepath.Join(
			"..",
			strings.Join(
				strings.Split(file, string(os.PathSeparator))[1:],
				string(os.PathSeparator),
			),
		)

		if err := os.MkdirAll(filepath.Dir(destPath), 0777); err != nil {
			warn(err)
			continue
		}

		if i, err := os.Stat(destPath); err == nil && i.IsDir() {
			warn("path '%s' is an existing directory, skipping")
			continue
		} else if *askFlag {
			fmt.Printf("file '%s' already exists, overwrite file? [n/Y] ", destPath)

			var ok string
			if _, err := fmt.Scanln(&ok); err != nil && err.Error() != "unexpected newline" {
				warn(err)
				continue
			}

			if !regexp.MustCompile("^[Yy].*$|^$").MatchString(ok) {
				continue
			}
		}

		// TODO os.O_RDWR -> os.O_WRONLY?
		destFile, err := os.OpenFile(destPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, info.Mode())
		if err != nil {
			warn(err)
			continue
		}
		defer destFile.Close()

		if err := tpl.ExecuteTemplate(destFile, tplName, ctx); err != nil {
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

func getXdgDir() string {
	xdgDir := os.Getenv("XDG_CONFIG_HOME")
	if xdgDir == "" {
		xdgDir = filepath.Join(os.Getenv("HOME"), ".config")
	}

	return filepath.Join(xdgDir, "templa")
}

func warn(a ...interface{}) {
	fmt.Fprintln(os.Stderr, programName+":", a)
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
