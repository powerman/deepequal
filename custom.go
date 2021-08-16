// Package deepequal provides improved reflect.DeepEqual.
//
// Differences from reflect.DeepEqual:
//
//   - If compared value implements `.Equal(valueOfSameType) bool` method then
//     it will be called instead of comparing values as is.
//   - If called `Equal` method will panics then whole DeepEqual will panics too.
//
// This means you can use this DeepEqual method to correctly compare types
// like time.Time or decimal.Decimal, without taking in account unimportant
// differences (like time zone or exponent).
package deepequal

import (
	"reflect"
	"unsafe"
)

// Disable check for unexported values.
func forceExported(v *reflect.Value) (undo func()) {
	ref := (*value)(unsafe.Pointer(v)) //nolint:gosec // Audit.
	flag := ref.flag
	ref.flag &^= flagRO
	return func() { ref.flag = flag }
}

func valueInterface(v reflect.Value) interface{} {
	undo := forceExported(&v)
	defer undo()
	return v.Interface()
}

func call(v reflect.Value, in []reflect.Value) []reflect.Value {
	undo := forceExported(&v)
	defer undo()
	for i := range in {
		undo := forceExported(&in[i])
		defer undo() //nolint:revive // By design.
	}
	return v.Call(in)
}

//nolint:gochecknoglobals // Const.
var boolType = reflect.TypeOf(true)

func equalFunc(v reflect.Value) (equal reflect.Value, ok bool) {
	equal = v.MethodByName("Equal")
	if !equal.IsValid() {
		return equal, false
	}
	typ := equal.Type()
	ok = typ.NumIn() == 1 && typ.In(0) == v.Type() &&
		typ.NumOut() == 1 && typ.Out(0) == boolType
	return equal, ok
}
