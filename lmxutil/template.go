package lmxutil

import (
	"encoding/xml"
	"io"

	"x-formation/lmx"
)

type license struct {
	XMLName  xml.Name `xml:LICENSEFILE`
	features []struct {
		XMLName  xml.Name `xml:FEATURE`
		name     string   `xml:NAME,attr`
		settings []struct {
			XMLName xml.Name `xml:SETTING`
		}
	}
}

func EncodeLicense(license interface{}, w io.Writer) error {
	return nil
}
