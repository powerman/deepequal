package deepequal

import (
	"reflect"
	"testing"
)

// This test should break in reflect.Value will be modified in
// incompatible way in future Go version.
func TestForceExported(t *testing.T) {
	t.Parallel()

	var v struct{ a int }
	val := reflect.ValueOf(v).Field(0)
	var res interface{}
	func() {
		defer func() { res = recover() }()
		_ = valueInterface(val)
	}()
	if res != nil {
		t.Errorf("valueInterface panics: %v", res)
	}
}
