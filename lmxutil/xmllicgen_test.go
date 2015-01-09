// +build ignore

package lmxutil

import (
	"testing"
)

func TestNewXmllicgen(t *testing.T) {
	x, err := NewXmllicgenFromPath()
	if err != nil {
		t.Fatalf(`expected err to be nil, got "%v" instead`, err)
	}
	if len(x.File()) == 0 {
		t.Error(`expected filename to be non-empty`)
	}
	if maj, _, _ := x.Version(); maj == 0 {
		t.Error(`expected major version != 0`)
	}
}
