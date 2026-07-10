// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE-go file.

package deepequal_test

import (
	"math"
	"reflect"
	"testing"

	"github.com/powerman/deepequal"
)

type Basic struct {
	x int
	y float32
}

type NotBasic Basic

type MyByte byte

type MyBytes []byte

type DeepEqualTest struct {
	a, b any
	eq   bool
}

// Simple functions for DeepEqual tests.
var (
	fn1 func()             // nil.
	fn2 func()             // nil.
	fn3 = func() { fn1() } // Not nil.
)

type self struct{}

type Loop *Loop

type Loopy any

var loop1, loop2 Loop

var loopy1, loopy2 Loopy

var cycleMap1, cycleMap2, cycleMap3 map[string]any

type structWithSelfPtr struct {
	p *structWithSelfPtr
	s string
}

func init() {
	loop1 = &loop2
	loop2 = &loop1

	loopy1 = &loopy2
	loopy2 = &loopy1

	cycleMap1 = map[string]any{}
	cycleMap1["cycle"] = cycleMap1
	cycleMap2 = map[string]any{}
	cycleMap2["cycle"] = cycleMap2
	cycleMap3 = map[string]any{}
	cycleMap3["different"] = cycleMap3
}

var deepEqualTests = []DeepEqualTest{
	// Equalities
	{nil, nil, true},
	{1, 1, true},
	{int32(1), int32(1), true},
	{0.5, 0.5, true},
	{float32(0.5), float32(0.5), true},
	{"hello", "hello", true},
	{make([]int, 10), make([]int, 10), true},
	{&[3]int{1, 2, 3}, &[3]int{1, 2, 3}, true},
	{Basic{1, 0.5}, Basic{1, 0.5}, true},
	{error(nil), error(nil), true},
	{map[int]string{1: "one", 2: "two"}, map[int]string{2: "two", 1: "one"}, true},
	{fn1, fn2, true},
	{[]byte{1, 2, 3}, []byte{1, 2, 3}, true},
	{[]MyByte{1, 2, 3}, []MyByte{1, 2, 3}, true},
	{MyBytes{1, 2, 3}, MyBytes{1, 2, 3}, true},

	// Inequalities
	{1, 2, false},
	{int32(1), int32(2), false},
	{0.5, 0.6, false},
	{float32(0.5), float32(0.6), false},
	{"hello", "hey", false},
	{make([]int, 10), make([]int, 11), false},
	{&[3]int{1, 2, 3}, &[3]int{1, 2, 4}, false},
	{Basic{1, 0.5}, Basic{1, 0.6}, false},
	{Basic{1, 0}, Basic{2, 0}, false},
	{map[int]string{1: "one", 3: "two"}, map[int]string{2: "two", 1: "one"}, false},
	{map[int]string{1: "one", 2: "txo"}, map[int]string{2: "two", 1: "one"}, false},
	{map[int]string{1: "one"}, map[int]string{2: "two", 1: "one"}, false},
	{map[int]string{2: "two", 1: "one"}, map[int]string{1: "one"}, false},
	{nil, 1, false},
	{1, nil, false},
	{fn1, fn3, false},
	{fn3, fn3, false},
	{[][]int{{1}}, [][]int{{2}}, false},
	{&structWithSelfPtr{p: &structWithSelfPtr{s: "a"}}, &structWithSelfPtr{p: &structWithSelfPtr{s: "b"}}, false},

	// Fun with floating point.
	{math.NaN(), math.NaN(), false},
	{&[1]float64{math.NaN()}, &[1]float64{math.NaN()}, false},
	{&[1]float64{math.NaN()}, self{}, true},
	{[]float64{math.NaN()}, []float64{math.NaN()}, false},
	{[]float64{math.NaN()}, self{}, true},
	{map[float64]float64{math.NaN(): 1}, map[float64]float64{1: 2}, false},
	{map[float64]float64{math.NaN(): 1}, self{}, true},

	// Nil vs empty: not the same.
	{[]int{}, []int(nil), false},
	{[]int{}, []int{}, true},
	{[]int(nil), []int(nil), true},
	{map[int]int{}, map[int]int(nil), false},
	{map[int]int{}, map[int]int{}, true},
	{map[int]int(nil), map[int]int(nil), true},

	// Mismatched types
	{1, 1.0, false},
	{int32(1), int64(1), false},
	{0.5, "hello", false},
	{[]int{1, 2, 3}, [3]int{1, 2, 3}, false},
	{&[3]any{1, 2, 4}, &[3]any{1, 2, "s"}, false},
	{Basic{1, 0.5}, NotBasic{1, 0.5}, false},
	{map[uint]string{1: "one", 2: "two"}, map[int]string{2: "two", 1: "one"}, false},
	{[]byte{1, 2, 3}, []MyByte{1, 2, 3}, false},
	{[]MyByte{1, 2, 3}, MyBytes{1, 2, 3}, false},
	{[]byte{1, 2, 3}, MyBytes{1, 2, 3}, false},

	// Possible loops.
	{&loop1, &loop1, true},
	{&loop1, &loop2, true},
	{&loopy1, &loopy1, true},
	{&loopy1, &loopy2, true},
	{&cycleMap1, &cycleMap2, true},
	{&cycleMap1, &cycleMap3, false},
}

func TestDeepEqual(t *testing.T) {
	for _, test := range deepEqualTests {
		if test.b == (self{}) {
			test.b = test.a
		}
		if r := deepequal.DeepEqual(test.a, test.b); r != test.eq {
			t.Errorf("DeepEqual(%#v, %#v) = %v, want %v", test.a, test.b, r, test.eq)
		}
	}
}

type Recursive struct {
	x int
	r *Recursive
}

func TestDeepEqualRecursiveStruct(t *testing.T) {
	a, b := new(Recursive), new(Recursive)
	*a = Recursive{12, a}
	*b = Recursive{12, b}
	if !deepequal.DeepEqual(a, b) {
		t.Error("DeepEqual(recursive same) = false, want true")
	}
}

type _Complex struct {
	a int
	b [3]*_Complex
	c *string
	d map[float64]float64
}

func TestDeepEqualComplexStruct(t *testing.T) {
	m := make(map[float64]float64)
	stra, strb := "hello", "hello"
	a, b := new(_Complex), new(_Complex)
	*a = _Complex{5, [3]*_Complex{a, b, a}, &stra, m}
	*b = _Complex{5, [3]*_Complex{b, a, a}, &strb, m}
	if !deepequal.DeepEqual(a, b) {
		t.Error("DeepEqual(complex same) = false, want true")
	}
}

func TestDeepEqualComplexStructInequality(t *testing.T) {
	m := make(map[float64]float64)
	stra, strb := "hello", "helloo" // Difference is here.
	a, b := new(_Complex), new(_Complex)
	*a = _Complex{5, [3]*_Complex{a, b, a}, &stra, m}
	*b = _Complex{5, [3]*_Complex{b, a, a}, &strb, m}
	if deepequal.DeepEqual(a, b) {
		t.Error("DeepEqual(complex different) = true, want false")
	}
}

type UnexpT struct {
	m map[int]int
}

func TestDeepEqualUnexportedMap(t *testing.T) {
	// Check that DeepEqual can look at unexported fields.
	x1 := UnexpT{map[int]int{1: 2}}
	x2 := UnexpT{map[int]int{1: 2}}
	if !deepequal.DeepEqual(&x1, &x2) {
		t.Error("DeepEqual(x1, x2) = false, want true")
	}

	y1 := UnexpT{map[int]int{2: 3}}
	if deepequal.DeepEqual(&x1, &y1) {
		t.Error("DeepEqual(x1, y1) = true, want false")
	}
}

var deepEqualPerfTests = []struct {
	x, y any
}{
	{x: int8(99), y: int8(99)},
	{x: []int8{99}, y: []int8{99}},
	{x: int16(99), y: int16(99)},
	{x: []int16{99}, y: []int16{99}},
	{x: int32(99), y: int32(99)},
	{x: []int32{99}, y: []int32{99}},
	{x: int64(99), y: int64(99)},
	{x: []int64{99}, y: []int64{99}},
	{x: int(999999), y: int(999999)},
	{x: []int{999999}, y: []int{999999}},

	{x: uint8(99), y: uint8(99)},
	{x: []uint8{99}, y: []uint8{99}},
	{x: uint16(99), y: uint16(99)},
	{x: []uint16{99}, y: []uint16{99}},
	{x: uint32(99), y: uint32(99)},
	{x: []uint32{99}, y: []uint32{99}},
	{x: uint64(99), y: uint64(99)},
	{x: []uint64{99}, y: []uint64{99}},
	{x: uint(999999), y: uint(999999)},
	{x: []uint{999999}, y: []uint{999999}},
	{x: uintptr(999999), y: uintptr(999999)},
	{x: []uintptr{999999}, y: []uintptr{999999}},

	{x: float32(1.414), y: float32(1.414)},
	{x: []float32{1.414}, y: []float32{1.414}},
	{x: float64(1.414), y: float64(1.414)},
	{x: []float64{1.414}, y: []float64{1.414}},

	{x: complex64(1.414), y: complex64(1.414)},
	{x: []complex64{1.414}, y: []complex64{1.414}},
	{x: complex128(1.414), y: complex128(1.414)},
	{x: []complex128{1.414}, y: []complex128{1.414}},

	{x: true, y: true},
	{x: []bool{true}, y: []bool{true}},

	{x: "abcdef", y: "abcdef"},
	{x: []string{"abcdef"}, y: []string{"abcdef"}},

	{x: []byte("abcdef"), y: []byte("abcdef")},
	{x: [][]byte{[]byte("abcdef")}, y: [][]byte{[]byte("abcdef")}},

	{x: [6]byte{'a', 'b', 'c', 'a', 'b', 'c'}, y: [6]byte{'a', 'b', 'c', 'a', 'b', 'c'}},
	{x: [][6]byte{{'a', 'b', 'c', 'a', 'b', 'c'}}, y: [][6]byte{{'a', 'b', 'c', 'a', 'b', 'c'}}},
}

func TestDeepEqualAllocs(t *testing.T) {
	for _, tt := range deepEqualPerfTests {
		t.Run(reflect.ValueOf(tt.x).Type().String(), func(t *testing.T) {
			got := testing.AllocsPerRun(100, func() {
				if !deepequal.DeepEqual(tt.x, tt.y) {
					t.Errorf("DeepEqual(%v, %v)=false", tt.x, tt.y)
				}
			})
			if int(got) != 0 {
				t.Errorf("DeepEqual(%v, %v) allocated %d times", tt.x, tt.y, int(got))
			}
		})
	}
}
