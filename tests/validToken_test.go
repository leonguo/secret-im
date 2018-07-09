package tests

import "testing"
import (
	"../app/lib"
	"strconv"
)

func TestValidToken(t *testing.T) {
	m, err := lib.DecodeToken("gKiLKriOxTMpcQ")
	if err != nil {
		t.Error("valid decode token err, got", err)
	}
	t.Log(m)
	t.Log(strconv.Quote(string(m[:])))
}
