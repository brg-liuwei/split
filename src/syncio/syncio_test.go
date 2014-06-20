package syncio

import (
	"os"
	"testing"
)

func TestSyncFile(t *testing.T) {
	err := os.MkdirAll("tmpPath", 0777)
	if err != nil {
		t.Error(err)
	}
	f1 := NewSyncFile("tmpPath/tmp")
	if f1 == nil {
		t.Error("open tmpPath/tmp err")
	}

	var pwd string
	pwd, err = os.Getwd()
	if err != nil {
		t.Error(err)
	}
	f2 := NewSyncFile(pwd + "/tmpPath/tmp")
	if f2 == nil {
		t.Error("open abs path error")
	}

	var n int
	var off int64

	s1 := "hello world"
	s2 := "Zia"

	n, err = f1.Write([]byte(s1), off)
	if err != nil {
		t.Error(err)
	}
	if n != len(s1) {
		t.Error("n = ", n, ", len(\"", s1, "\") = ", len(s1))
	}

	off += int64(n)

	n, err = f2.Write([]byte(s2), off)
	if err != nil {
		t.Error(err)
	}
	if n != len(s2) {
		t.Error("n = ", n, ", len(\"", s1, "\") = ", len(s1))
	}

	off += int64(n)

	buf := make([]byte, 100)
	n, err = f1.Read(buf, 0)
	if err != nil {
		t.Error(err)
	}
	if int64(n) != off {
		t.Error("n = ", n, ", off = ", off)
	}
	if n != len(s1+s2) {
		t.Error("n = ", n, ", len(", s1+s2, ") = ", len(s1+s2))
	}
	buf = buf[:n]
	if string(buf) != s1+s2 {
		t.Error("string(buf): ", string(buf), ", s1 + s2: ", s1+s2)
	}

	err = f1.Close()
	if err != nil {
		t.Error("close f1 error: ", err)
	}
	err = f2.Close()
	if err != nil {
		t.Error("close f2 error: ", err)
	}

	err = os.RemoveAll("tmpPath")
	if err != nil {
		t.Error(err)
	}
}
