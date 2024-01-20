package maps

import (
	"errors"
	"reflect"
)

func Merge[M1, M2 ~map[K]V, K comparable, V any](dst M1, src M2) {
	if dst == nil {
		return
	}

	merge(reflect.ValueOf(dst), reflect.ValueOf(src))
}

var errDifferentTypes = errors.New("cannot merge non-assignable types")
var mergeFunc reflect.Value

func init() {
	mergeFunc = reflect.ValueOf(merge)
}

func merge(dst, src reflect.Value) {
	it := src.MapRange()
	for it.Next() {
		sK := it.Key()
		sV := it.Value()

		if dV := dst.MapIndex(sK); dV.IsValid() {
			for dV.Kind() == reflect.Interface || dV.Kind() == reflect.Pointer {
				dV = dV.Elem()
			}
			for sV.Kind() == reflect.Interface || sV.Kind() == reflect.Pointer {
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
					v := dV
					for i := 0; i < sV.Len(); i++ {
						v = reflect.Append(v, sV.Index(i))
					}
					sV = v
				} else {
					panic(errDifferentTypes)
				}
			default:
				if !sV.Type().AssignableTo(dV.Type()) {
					panic(errDifferentTypes)
				}
			}

			if dV.CanSet() {
				dV.Set(sV)
				continue
			}
		}

		dst.SetMapIndex(sK, sV)
	}
}
