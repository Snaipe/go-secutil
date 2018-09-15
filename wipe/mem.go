/* go-secutil
 *
 * Copyright (C) 2018  Franklin "Snaipe" Mathieu <me@snai.pe>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package wipe

import (
	"fmt"
	"reflect"
	"unsafe"
)

// ZeroUnsafe fills the memory in the range [addr,addr+size) with zeroes.
func ZeroUnsafe(addr unsafe.Pointer, size uintptr) {
	var buf []byte
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
	hdr.Data = uintptr(addr)
	hdr.Len = int(size)
	hdr.Cap = hdr.Len

	for i := range buf {
		buf[i] = 0
	}
}

// Zero wipes the memory occupied by v by filling it with zeroes.
//
// If v is a slice, the function zeroes out v[0:cap(v)], preserving the length
// and the capacity of the slice.
//
// If v is a string, the function zeroes out v[0:len(v)], preserving the length
// of the string. Note that since the string is mutated, and strings in go are
// read-only constructs, the caller must make sure that the underlying memory
// range is writeable.
//
// If v is a pointer, the function wipes the memory in the range
// [v, v+unsafe.Sizeof(*v)).
//
// It panics if v is neither a pointer, a slice, or a string.
func Zero(v interface{}) {
	ZeroValue(reflect.ValueOf(v))
}

// ZeroValue behaves like Zero, but acts on a reflect.Value instead.
func ZeroValue(v reflect.Value) {
	var addr, size uintptr
	switch v.Kind() {
	case reflect.String:
		/* we have to use this hack to access the underlying pointer to the
		   string data -- this works because value assignation for strings
		   does not copy the data, but only the string header. That way,
		   we can copy the header here and access it directly */
		str := v.String()
		hdr := (*reflect.StringHeader)(unsafe.Pointer(&str))
		addr = hdr.Data
		size = uintptr(hdr.Len)
	case reflect.Slice:
		addr = v.Pointer()
		size = uintptr(v.Cap()) * v.Type().Elem().Size()
	case reflect.Ptr:
		v = v.Elem()
		addr = v.UnsafeAddr()
		size = v.Type().Size()
	default:
		panic(fmt.Sprintf("%v is neither a pointer, a slice, or a string.", v.Type().Name()))
	}

	ZeroUnsafe(unsafe.Pointer(addr), size)
}

func isCompound(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Interface, reflect.Ptr, reflect.Struct, reflect.Map, reflect.Array, reflect.Slice:
		return true
	}
	return false
}

func nestedWipe(v reflect.Value) {
	switch v.Kind() {
	case reflect.Interface, reflect.Ptr:
		ZeroValueDeep(v.Elem())
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			nestedWipe(v.Field(i))
		}
	case reflect.Array, reflect.Slice:
		if isCompound(v.Type().Elem()) {
			for i := 0; i < v.Len(); i++ {
				nestedWipe(v.Index(i).Elem())
			}
		}
	case reflect.Map:
		vtype := v.Type().Elem()
		if isCompound(vtype) {
			for _, k := range v.MapKeys() {
				ZeroValueDeep(v.MapIndex(k))
			}
		}
		zero := reflect.Zero(vtype)
		for _, k := range v.MapKeys() {
			v.SetMapIndex(k, zero)
		}
	}
}

// pointers and indexables can be made "addressable" by accessing their
// data pointers
func canAddr(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.String:
		return true
	}
	return false
}

// ZeroDeep recursively zeroes out the contents of v.
//
// Maps at any depth are semi-unsupported: since map values are unaddressable,
// they can't be cleared reliably, unless the type of the value isn't itself
// sensitive, but holds a reference to sensitive data, i.e. pointers, slices,
// and strings.
func ZeroDeep(v interface{}) {
	ZeroValueDeep(reflect.ValueOf(v))
}

// ZeroValueDeep behaves like ZeroDeep, but acts on a reflect.Value instead.
func ZeroValueDeep(v reflect.Value) {
	nestedWipe(v)
	if canAddr(v.Type()) && v.CanAddr() {
		ZeroValue(v)
	}
}
