// +build ignore

package lmxutil

import (
	"errors"
	"os/exec"
	"regexp"
	"strconv"
)

type Xmllicgen interface {
	File() string
	Version() (int, int, int)
	Generate(...interface{}) ([]string, error)
	GenerateTo([]string, ...interface{}) error
}

type xmllicgen struct {
	file    string
	version [3]int
}

var re = regexp.MustCompile(`LM-X XML License Generator v(\d)\.(\d)\.{0,1}(\d){0,1}`)

// ErrVersionExtract
var ErrVersionExtract = errors.New(`lmxutil: unable to extract version information`)

func NewXmllicgen(file string) (Xmllicgen, error) {
	x := &xmllicgen{file: file}
	blob, err := exec.Command(file).CombinedOutput()
	if _, exit := err.(*exec.ExitError); err != nil && !exit {
		return nil, err
	}
	match := re.FindStringSubmatch(string(blob))
	if match == nil || len(match) == 0 {
		return nil, ErrVersionExtract
	}
	match = match[1:]
	if len(match) != 2 && len(match) != 3 {
		return nil, ErrVersionExtract
	}
	for i := range match {
		if x.version[i], err = strconv.Atoi(match[i]); err != nil {
			return nil, ErrVersionExtract
		}
	}
	return x, nil
}

func NewXmllicgenFromPath() (Xmllicgen, error) {
	file, err := exec.LookPath("xmllicgen")
	if err != nil {
		return nil, err
	}
	return NewXmllicgen(file)
}

func (x *xmllicgen) File() string {
	return x.file
}

func (x *xmllicgen) Version() (int, int, int) {
	return x.version[0], x.version[1], x.version[2]
}

func (x *xmllicgen) Generate(lic ...interface{}) ([]string, error) {
	return nil, nil
}

func (x *xmllicgen) GenerateTo(output []string, lic ...interface{}) error {
	return nil
}
