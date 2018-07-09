package util

import "regexp"

// 验证手机号码格式
func IsValidNumber(phone string) bool {
	if check, _ := regexp.MatchString("^([+]?\\d{1,4}[-\\s]?|)\\d{3}[-\\s]?\\d{3}[-\\s]?\\d{4}", phone); !check {
		return false
	}
	return true
}