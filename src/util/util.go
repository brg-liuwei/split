package util

import (
	"io"
)

var vintHead int64 = 0x80
var vintMask int64 = 0x7F

func ReadVint(buf []byte) (val int64, n int, err error) {
	for i, shift := 0, uint(0); i < len(buf); i, shift = i+1, shift+7 {
		b := int64(buf[i])
		val |= (b & vintMask) << shift
		if b&vintHead == 0 {
			n = i + 1
			return
		}
	}
	val = 0
	err = io.EOF
	return
}

func WriteVint(val int64, buf []byte) (n int, err error) {
	for i := 0; i < len(buf); i++ {
		if i == 9 { /* 7 * 9 = 63, 此时val的最高位为1 */
			buf[i] = 0x1
			n = i + 1
			return
		}
		b := byte(val & 0x7F)
		buf[i] = b | 0x80
		val >>= 7
		if val == 0 {
			buf[i] &= 0x7F
			n = i + 1
			return
		}
	}
	err = io.EOF
	return
}

func ReadVstr(buf []byte) (str string, n int, err error) {
	var v int64
	v, n, err = ReadVint(buf)
	if err != nil {
		return
	}
	buf = buf[n:]
	strlen := int(v)
	if strlen > len(buf) {
		n = 0
		err = io.EOF
		return
	}
	n += strlen
	str = string(buf[:strlen])
	return
}

func WriteVstr(str string, buf []byte) (n int, err error) {
	strlen := len(str)
	n, err = WriteVint(int64(strlen), buf)
	if err != nil {
		return
	}
	buf = buf[n:]
	if strlen > len(buf) {
		n = 0
		err = io.EOF
		return
	}
	copy(buf[:], []byte(str))
	n += strlen
	return
}

func ReadVbytes(buf []byte) (b []byte, n int, err error) {
	var v int64
	v, n, err = ReadVint(buf)
	if err != nil {
		return
	}
	buf = buf[n:]
	blen := int(v)
	if len(buf) >= blen {
		n += blen
		b = buf[:blen]
		return
	}
	err = io.EOF
	return
}

func WriteVbytes(src []byte, dst []byte) (n int, err error) {
	blen := len(src)
	n, err = WriteVint(int64(blen), dst)
	if err != nil {
		return
	}
	dst = dst[n:]
	if len(dst) < len(src) {
		n = 0
		err = io.EOF
		return
	}
	n += len(src)
	copy(dst, src)
	return
}
