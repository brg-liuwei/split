package util

import (
	"io"
	"testing"
)

func TestVint(t *testing.T) {
	var v1 int64 = 100
	var v2 int64 = 0
	var v3 int64 = -1
	var v4 int64 = 1024 * 1024 * 1024

	testReadWriteVint(t, v1)
	testReadWriteVint(t, v2)
	testReadWriteVint(t, v3)
	testReadWriteVint(t, v4)
}

func testReadWriteVint(t *testing.T, v int64) {

	buf := make([]byte, 20)

	n, err := WriteVint(v, buf)
	if err != nil {
		t.Error(err)
	}

	val, n1, e := ReadVint(buf)
	if e != nil {
		t.Error(err)
	}
	if n != n1 {
		t.Error("n: ", n, ", n1: ", n1)
	}
	if val != v {
		t.Error("val: ", val, "v: ", v)
	}
}

func TestEOF(t *testing.T) {
	buf := make([]byte, 1)
	_, err := WriteVint(1024, buf)
	if err != io.EOF {
		t.Error(err, ", Should io.EOF")
	}

	var buf1 []byte = []byte{0xDE, 0xAD, 0xBE, 0xEF}
	_, _, e := ReadVint(buf1)
	if e != io.EOF {
		t.Error(e, ", Should io.EOF")
	}
}

func TestVstr(t *testing.T) {
	buf := make([]byte, 12)
	n, err := WriteVstr("hello world", buf)
	if err != nil {
		t.Error(err)
	}
	println("n = ", n)
	s, n1, e := ReadVstr(buf)
	if e != nil {
		t.Error(e)
	}
	if n1 != n {
		t.Error("n1: ", n1, ", n: ", n)
	}
	if s != "hello world" {
		t.Error("s: ", s)
	}
}

func TestVbytes(t *testing.T) {
	orig := []byte("hello world")
	dst := make([]byte, 12)
	n, err := WriteVbytes(orig, dst)
	if err != nil {
		t.Error(err)
	}
	if n != 12 {
		t.Error("n = ", n, ", n should be 12")
	}
	result, n1, e := ReadVbytes(dst)
	if e != nil {
		t.Error(e)
	}
	if n1 != 12 {
		t.Error("n = ", n, ", n should be 12")
	}
	if len(result) != len(orig) {
		t.Error("len(result) = ", len(result), ", len(orig) = ", len(orig))
	}
	for i, v := range result {
		if orig[i] != v {
			t.Error("i: ", i, ", v: ", v, ", orig[i]: ", orig[i])
		}
	}
}
