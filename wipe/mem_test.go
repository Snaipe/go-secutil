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
	"testing"
)

// probably needs some improvements, but good enough for now

func TestZeroSlice(t *testing.T) {
	buf := []byte{0xdb, 0xdb}

	buf = buf[:1]
	Zero(buf)
	buf = buf[:cap(buf)]

	if buf[0] != 0 || buf[1] != 0 {
		t.Fatal("expected slice to be wiped")
	}
	if len(buf) != 2 {
		t.Fatal("expected slice to not be niled")
	}
}

func TestZeroSliceAddr(t *testing.T) {
	buf := []byte{0xdb, 0xdb}
	same := buf

	buf = buf[:1]
	Zero(&buf)

	if same[0] != 0xdb || same[1] != 0xdb {
		t.Fatal("expected slice to not be wiped")
	}
	if buf != nil {
		t.Fatal("expected slice to be niled")
	}
}

func TestZeroArray(t *testing.T) {
	arr := [...]byte{0xdb, 0xdb}

	Zero(&arr)

	if arr[0] != 0 || arr[1] != 0 {
		t.Fatal("expected array be wiped")
	}
}

func TestZeroString(t *testing.T) {
	s := string([]byte{'a', 'b'})

	Zero(s)

	if s[0] != 0 || s[1] != 0 {
		t.Fatal("expected string be wiped")
	}
}

func TestZeroStruct(t *testing.T) {
	s := struct {
		field1 uint32
		field2 uint32
	}{ 0xdbdbdbdb, 0xdbdbdbdb }

	Zero(&s)

	if s.field1 != 0 || s.field2 != 0 {
		t.Fatal("expected struct to be wiped")
	}
}
