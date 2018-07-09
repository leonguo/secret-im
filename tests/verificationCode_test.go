package tests

import (
	"testing"
	"../app/lib"
	"../config"
)

func TestValidCode(t *testing.T) {
	config.Init()
    m := lib.VerificationCode("+8618575682804")
	if m != "123456" {
		t.Error("valid phone 8618575682804 code err, got", m)
	}

	m1 := lib.VerificationCode("+8618575682801")
	if m1 == "123456" {
		t.Error("valid phone 8618575682804 code err, got", m1)
	}
}