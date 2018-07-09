package tests

import (
	"testing"
	"encoding/base64"
	"log"
	"fmt"
)

func TestValidString(t *testing.T) {
	input := []byte("dss")
	encodeString := base64.StdEncoding.EncodeToString(input)
	fmt.Println(encodeString)
	str := "Mwj7za8HEiEF8bJLYn8Wq/n1GJB642yp6uoItHrX+wbmYwOI04awHVAaIQWXJBWRxtTY+sWbDD/I2+13YfkxORo6qfm6DyokUCShSyLTATMKIQUEFar76wOt8nZJJ/YsAn1ekR5fRd8SPfMTohoO7t/5fxABGAAioAHV01ji06sjB2q1xlxDs8N4+rBPhTsXHfGlmSkCjCHVq7EMSTijPCOgygYAQnMG3gNKtyO68rxAswB5srM8zAlEagJctz14rd4uJyit1R+TbfmfBgCROdeL/p3Kko007zY8HwnofuzuwqNsa2C6n5WMq/eK4127nEDop9efTsHoeNU04oHC2MfNcOtgjPn0QZhbMDwlooUXx+sHw1YGRoDiQlgjp8cTvGcokVMwkLvnAw=="
	content, err := base64.StdEncoding.DecodeString(str)
	log.Printf("content %v", string(content))
	log.Printf("content %v", err)
	if "" == "123456" {
		t.Error("valid phone 8618575682804 code err, got")
	}
}
