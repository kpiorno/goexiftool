package dtos

import "encoding/xml"

type ExiftoolXML struct {
	TagInfo string  `xml:"taginfo" json:"taginfo"`
	Tables  []Table `xml:"table" json:"table"`
}

type Table struct {
	XMLName xml.Name `xml:"table" json:"-"`
	Name    string   `xml:"name,attr" json:"name"`
	Tag     Tag      `xml:"tag" json:"tag"`
}

type Tag struct {
	XMLName  xml.Name `xml:"tag" json:"-"`
	Name     string   `xml:"name,attr" json:"name"`
	Writable bool     `xml:"writable,attr" json:"writable"`
	Desc     []Desc   `xml:"desc" json:"desc"`
}

type Desc struct {
	Value string `xml:",chardata" json:"value"`
	Lang  string `xml:"lang,attr" json:"lang"`
}
