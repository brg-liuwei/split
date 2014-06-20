package main

import (
	"config"
	"fmt"
	"kv"
	"syncio"
	"util"
)

var f syncio.SyncFile

func main() {
	xml := config.NewXmlDecoder("conf/schema.xml")
	for xml.HasNext() {
		name, attr := xml.Next()

		// do some stuff
		_ = name
		_ = attr
	}

	buf := make([]byte, 20)
	util.WriteVint(int64(-1), buf)
	v, _, _ := util.ReadVint(buf)
	fmt.Println("v = ", v)

	name := "testDb"
	kv.Open(&name)
	k := "hello"
	val := "world"
	kv.Put(&k, &val)
	v1, _ := kv.Get(&k)
	println("v1 = ", v1)
	kv.Close()
	_ = f
}
