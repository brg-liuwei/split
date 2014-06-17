package config

import (
	"testing"
)

var testXmlDecoder *XmlDecoder

func TestNewXmlDecoder(t *testing.T) {
	testXmlDecoder = NewXmlDecoder("schema_test.xml")
	if testXmlDecoder == nil {
		t.Error("Test New Xml Decoder Error")
	}
}

func TestReadXml(t *testing.T) {
	p := testXmlDecoder
	if !p.HasNext() {
		t.Error("Cannot get schema")
	}

	// schema
	name, m := p.Next()
	if name != "schema" {
		t.Error()
	}
	if len(m) != 2 {
		t.Error()
	}
	if m["name"] != "usr-tag" {
		t.Error()
	}

	// fields
	if !p.HasNext() {
		t.Error()
	}
	name, m = p.Next()
	if name != "fields" {
		t.Error()
	}
	if len(m) != 0 {
		t.Error()
	}

	// field uid
	if !p.HasNext() {
		t.Error()
	}
	name, m = p.Next()
	if name != "field" {
		t.Error()
	}
	if len(m) != 6 {
		t.Error()
	}
	if m["name"] != "uid" || m["type"] != "str" || m["length"] != "0" || m["indexed"] != "false" || m["stored"] != "true" || m["filter"] != "false" {
		t.Error()
	}

	// prime
	if !p.HasNext() {
		t.Error()
	}
	name, m = p.Next()
	if name != "prime" {
		t.Error()
	}
	if len(m) != 1 {
		t.Error()
	}
	if m["key"] != "uid" {
		t.Error()
	}

	// end
	if p.HasNext() {
		t.Error()
	}
}

func TestReStart(t *testing.T) {
	testXmlDecoder.ReStart()
	TestReadXml(t)
}
