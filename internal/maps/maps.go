package maps

import (
	"errors"
	"reflect"
)

var (
	ErrDifferentTypes = errors.New("cannot merge non-assignable types")

	mergeFunc reflect.Value
)

func init() {
	mergeFunc = reflect.ValueOf(merge)
}

func Merge[M1, M2 ~map[K]V, K comparable, V any](dst M1, src M2) error {
	if dst == nil {
		return nil
	}

	return merge(reflect.ValueOf(dst), reflect.ValueOf(src))
}

func merge(dst, src reflect.Value) error {
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
					if r := mergeFunc.Call(args); r != nil {
						err, _ := r[0].Interface().(error)
						return err
					}
					continue
				} else {
					return ErrDifferentTypes
				}
			case reflect.Slice:
				if sV.Kind() == reflect.Slice && sV.Type().Elem().AssignableTo(dV.Type().Elem()) {
					for i := 0; i < sV.Len(); i++ {
						dV = reflect.Append(dV, sV.Index(i))
					}
					sV = dV
				} else {
					return ErrDifferentTypes
				}
			default:
				if !sV.Type().AssignableTo(dV.Type()) {
					return ErrDifferentTypes
				}
			}
		}

		dst.SetMapIndex(sK, sV)
	}

	return nil
}
