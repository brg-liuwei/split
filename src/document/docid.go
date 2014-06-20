package document

import (
	"errors"
	"fmt"
	"kv"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"syncio"
)

var curId int64 = 0
var appender *os.File
var delHandler *syncio.SyncFile
var once sync.Once

const dbname string = "uid-db"
const delfile string = "_docid.del"

func DocInit() {
	once.Do(func() {
		var err error
		name := dbname
		kv.Open(&name)
		appender, err = os.OpenFile(delfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		delHandler = syncio.NewSyncFile(delfile)
		if delHandler == nil {
			panic("open SyncFile: " + delfile + " error")
		}
	})
}

func AllocId() int64 {
	id := atomic.AddInt64(&curId, 1) - 1
	if id&0x7 == 0 {
		var b [1]byte
		_, err := appender.Write(b[:])
		if err != nil {
			panic(err)
		}
	}
	return id
}

func readIdByte(id int64) (byte, error) {
	var buf [1]byte
	off := id >> 3
	n, err := delHandler.Read(buf[:], off)
	if err != nil {
		fmt.Println("pread docid.del error: ", err)
		return byte(0), err
	}
	if n != 1 {
		fmt.Println("error: pread docid.del return: ", n)
		return byte(0), errors.New("pread docid.del return " + strconv.Itoa(n))
	}
	return buf[0], nil
}

func writeIdByte(id int64, b byte) (int, error) {
	var buf [1]byte
	buf[0] = b
	off := id >> 3
	return delHandler.Write(buf[:], off)
}

func ActiveId(id int64) bool {
	if id >= curId {
		return false
	}
	bit := uint(id & 0x7)
	b, err := readIdByte(id)
	if err != nil {
		fmt.Println(err)
		return false
	}
	mask := byte(0x80) >> bit
	if b&mask != 0 {
		return false
	}
	return true
}

func DeleteId(id int64) {
	if id >= curId {
		return
	}
	b, err := readIdByte(id)
	if err != nil {
		fmt.Println(err)
		return
	}
	bit := uint(id & 0x7)
	mask := byte(0x80) >> bit
	b |= mask
	writeIdByte(id, b)
}
