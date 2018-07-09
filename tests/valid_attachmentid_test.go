package tests

import "testing"
import (
	"../util"
)

func TestAttachmentID(t *testing.T) {
	m, err := util.GenerateAttachmentId()
	if err != nil {
		t.Error("gen attachment error ", m)
	}
}