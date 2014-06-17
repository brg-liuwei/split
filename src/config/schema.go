package config

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
)

type XmlDecoder struct {
	c    []byte
	d    *xml.Decoder
	name string
	attr map[string]string
}

func NewXmlDecoder(path string) *XmlDecoder {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("NewXmlDecoder Error: ", err)
		return nil
	}
	d := xml.NewDecoder(bytes.NewBuffer(content))
	return &XmlDecoder{c: content, d: d}
}

func (x *XmlDecoder) ReStart() {
	x.d = xml.NewDecoder(bytes.NewBuffer(x.c))
	x.name = ""
	x.attr = nil
}

func (x *XmlDecoder) HasNext() bool {
	for {
		token, err := x.d.Token()
		if err != nil {
			if err == io.EOF {
				return false
			}
			fmt.Println("Error in HasNext: ", err)
			return false
		}
		switch t := token.(type) {
		case xml.StartElement:
			x.name = t.Name.Local
			x.attr = make(map[string]string)
			for _, attr := range t.Attr {
				x.attr[attr.Name.Local] = attr.Value
			}
			return true
		case xml.EndElement:
		case xml.CharData:
		}
	}
}

func (x *XmlDecoder) Next() (name string, attr map[string]string) {
	name = x.name
	attr = x.attr
	return
}
