package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type stringSlice []string

func (sl *stringSlice) String() string {
	return fmt.Sprint(*sl)
}

func (sl *stringSlice) Set(value string) error {
	*sl = append(*sl, value)
	return nil
}

func warn(a ...interface{}) {
	fmt.Fprintln(os.Stderr, programName+":", a)
}

func getXdgDir() string {
	xdgDir := os.Getenv("XDG_CONFIG_HOME")
	if xdgDir == "" {
		xdgDir = filepath.Join(os.Getenv("HOME"), ".config")
	}

	return filepath.Join(xdgDir, "templa")
}
