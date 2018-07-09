package tests

import "testing"
import (
	"../util"
	"fmt"
)

func TestAES(t *testing.T) {
	contentByte := []byte{33, 33, 4, 5, 4, 6, 7, 7, 7, 8, 8, 8, 9, 9, 1}
	m, err := util.AESCBCEncrypt(contentByte, "9n4mcqBAUhUdLO0RuzkXJI1rI2Nz5h5cMUN3+zFMQYpw+2Aq/Jn60WVeUU0XcoX0b9EAtA==")
	if err != nil {
		t.Error("valid 8618575682804 err, got", err)
	}
	fmt.Printf("mmm >>>%v", m)
}
