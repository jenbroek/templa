package maps

import (
	"errors"
	"reflect"
)

var (
	errDifferentTypes = errors.New("cannot merge non-assignable types")

	mergeFunc reflect.Value
)

func init() {
	mergeFunc = reflect.ValueOf(merge)
}

func Merge[M1, M2 ~map[K]V, K comparable, V any](dst M1, src M2) {
	if dst == nil {
		return
	}

	merge(reflect.ValueOf(dst), reflect.ValueOf(src))
}

func merge(dst, src reflect.Value) {
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
					mergeFunc.Call(args)
					continue
				} else {
					panic(errDifferentTypes)
				}
			case reflect.Slice:
				if sV.Kind() == reflect.Slice && sV.Type().Elem().AssignableTo(dV.Type().Elem()) {
					for i := 0; i < sV.Len(); i++ {
						dV = reflect.Append(dV, sV.Index(i))
					}
					sV = dV
				} else {
					panic(errDifferentTypes)
				}
			default:
				if !sV.Type().AssignableTo(dV.Type()) {
					panic(errDifferentTypes)
				}
			}
		}

		dst.SetMapIndex(sK, sV)
	}
}
