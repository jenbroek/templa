package main

import "fmt"

var errDifferentTypes = fmt.Errorf("cannot merge different map/slice types")

func mergeMaps(dst map[string]any, src map[string]any) error {
	if dst == nil {
		return nil
	}

	for sK, sV := range src {
		if dV, ok := dst[sK]; ok {
			switch dV := dV.(type) {
			case map[string]any:
				if v, ok := sV.(map[string]any); ok {
					if err := mergeMaps(dV, v); err != nil {
						return err
					}
					continue
				} else {
					return errDifferentTypes
				}
			case []any:
				if v, ok := sV.([]any); ok {
					sV = append(dV, v...)
				} else {
					return errDifferentTypes
				}
			}
		}

		dst[sK] = sV
	}

	return nil
}
