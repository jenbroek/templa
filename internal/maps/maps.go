package maps

import "reflect"

type MergeError struct {
	SrcType string
	DstType string
}

func (e *MergeError) Error() string {
	return "maps: cannot assign type " + e.SrcType + " to type " + e.DstType
}

func Merge[M1, M2 ~map[K]V, K comparable, V any](dst M1, src M2) {
	if dst == nil {
		return
	}

	merge(reflect.ValueOf(dst), reflect.ValueOf(src))
}

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
				if assignableTo(sV, dV, reflect.Map) {
					args := []reflect.Value{reflect.ValueOf(dV), reflect.ValueOf(sV)}
					mergeFunc.Call(args)
					continue
				}
			case reflect.Slice:
				if assignableTo(sV, dV, reflect.Slice) {
					v := dV
					for i := 0; i < sV.Len(); i++ {
						v = reflect.Append(v, sV.Index(i))
					}
					sV = v
				}
			default:
				if !sV.Type().AssignableTo(dV.Type()) {
					panic(&MergeError{sV.Type().String(), dV.Type().String()})
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

func assignableTo(v1, v2 reflect.Value, kind reflect.Kind) bool {
	if v1.Kind() == kind {
		v1ElemType := v1.Type().Elem()
		v2ElemType := v2.Type().Elem()

		if v1ElemType.AssignableTo(v2ElemType) {
			return true
		}
	}

	panic(&MergeError{v1.Type().String(), v2.Type().String()})
}
