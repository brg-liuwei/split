package main

import (
	"config"
	"fmt"
	"util"
)

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
}
