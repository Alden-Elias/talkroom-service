package myUtils

import (
	"regexp"
)

func IsAnEmail(email string) bool {
	res, _ := regexp.MatchString("^[a-zA-Z0-9_.-]+@[a-zA-Z0-9-]+(\\.[a-zA-Z0-9-]+)*\\.[a-zA-Z0-9]{2,6}$", email)
	return res
}

func PasswordCheck(passwd string) bool {
	res, _ := regexp.MatchString("^[\\da-zA-Z-]{6,20}$", passwd)
	return res
}
