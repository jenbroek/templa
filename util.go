package main

import (
	"fmt"
	"io/fs"
	"os"
	"reflect"
)

// See: https://github.com/golang/go/issues/44286
type osFS struct{}

func (*osFS) Open(name string) (fs.File, error) { return os.Open(name) }

var (
	errDifferentTypes = fmt.Errorf("cannot merge non-assignable types")

	mergeMapsInternalFunc reflect.Value
)

func init() {
	mergeMapsInternalFunc = reflect.ValueOf(mergeMapsInternal)
}

func MergeMaps[M1, M2 ~map[K]V, K comparable, V any](dst M1, src M2) error {
	if dst == nil {
		return nil
	}

	return mergeMapsInternal(reflect.ValueOf(dst), reflect.ValueOf(src))
}

func mergeMapsInternal(dst, src reflect.Value) error {
	it := src.MapRange()
	for it.Next() {
		sK := it.Key()
		sV := it.Value()

		if dV := dst.MapIndex(sK); dV.IsValid() {
			if dV.Kind() == reflect.Interface {
				dV = dV.Elem()
			}
			if sV.Kind() == reflect.Interface {
				sV = sV.Elem()
			}

			switch dV.Kind() {
			case reflect.Map:
				if sV.Kind() == reflect.Map && sV.Type().Elem().AssignableTo(dV.Type().Elem()) {
					args := []reflect.Value{reflect.ValueOf(dV), reflect.ValueOf(sV)}
					if r := mergeMapsInternalFunc.Call(args); r != nil {
						err, _ := r[0].Interface().(error)
						return err
					}
					continue
				} else {
					return errDifferentTypes
				}
			case reflect.Slice:
				if sV.Kind() == reflect.Slice && sV.Type().Elem().AssignableTo(dV.Type().Elem()) {
					for i := 0; i < sV.Len(); i++ {
						dV = reflect.Append(dV, sV.Index(i))
					}
					sV = dV
				} else {
					return errDifferentTypes
				}
			default:
				if !sV.Type().AssignableTo(dV.Type()) {
					return errDifferentTypes
				}
			}
		}

		dst.SetMapIndex(sK, sV)
	}

	return nil
}
