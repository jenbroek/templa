package main

import (
	"os"
	"path/filepath"
)

func mergeMaps(out map[string]any, in map[string]any) {
	for k, v := range in {
		if existingVal, ok := out[k]; ok {
			switch existingVal := existingVal.(type) {
			case map[string]any:
				if nv, ok := v.(map[string]any); ok {
					mergeMaps(existingVal, nv)
					continue
				}
			case []any:
				if nv, ok := v.([]any); ok {
					v = append(existingVal, nv...)
				}
			}
		}

		out[k] = v
	}
}

func relPath(tmplPath string) (string, error) {
	if !filepath.IsAbs(tmplPath) {
		return tmplPath, nil
	}

	invokePath, err := os.Executable()
	if err != nil {
		return "", err
	}

	return filepath.Rel(filepath.Dir(invokePath), tmplPath)
}
