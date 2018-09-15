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

package secutil

import (
	"reflect"
	"unsafe"
)

// BytesToStringView returns a string backed by the same underlying memory
// as buf.
//
// In Go, converting byte slices to strings always copy the contents of the
// slice, which can sometimes be unfortunate when converting a byte slice
// containing sensitive data to a string, as such a slice is usually allocated
// by the user outside of the garbage collector, and the resulting string
// ends up being managed by the GC.
func BytesToStringView(buf []byte) (str string) {
	bhdr := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
	shdr := (*reflect.StringHeader)(unsafe.Pointer(&str))
	shdr.Data = bhdr.Data
	shdr.Len = bhdr.Len
	return
}

// StringToBytesView returns a byte slice backed by the same underlying memory
// as str.
//
// Note that this can be dangerous if the underlying memory is read-only, as
// byte slices are assumed to be mutable. This can happen when the string comes
// from a string literal, and users of this function are expected to know what
// they are doing with the resulting slice.
func StringToBytesView(str string) (buf []byte) {
	bhdr := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
	shdr := (*reflect.StringHeader)(unsafe.Pointer(&str))
	bhdr.Data = shdr.Data
	bhdr.Len = shdr.Len
	bhdr.Cap = shdr.Len
	return
}
