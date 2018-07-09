package tests

import "testing"
import (
	"../app/lib"
	"bytes"
)

func TestValidContactToken(t *testing.T) {
	hs := lib.ContactToken("+8618575682804")
	if hs == nil {
		t.Error("valid decode token err, got", hs)
	}
	var b bytes.Buffer
	b.Write(hs)
	b.Truncate(10)
	hs = b.Bytes()
	t.Log(hs)
	t.Log(string(hs))
}
