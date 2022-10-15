package dao

import (
	mathRand "math/rand"
	"strconv"
	"talkRoom/myUtils"
	"time"
)

//GetVerificationCode 在Redis中保存验证码
func GetVerificationCode(email string) (vCode string, err error) {
	if exist, err1 := redisDb.Exists("Register:Verification:" + email).Result(); err1 != nil {
		err = err1
	} else if exist == 0 {
		vCode = strconv.Itoa(mathRand.Intn(9000) + 1000)
		err = redisDb.Set("Register:Verification:"+email, vCode, 5*time.Minute).Err()
	} else {
		vCode, err = redisDb.Get("Register:Verification:" + email).Result()
	}
	return
}

func VerifyCode(email, vCode string) (isMatch bool, err error) {
	if exist, err1 := redisDb.Exists("Register:Verification:" + email).Result(); err1 != nil {
		isMatch, err = false, err1
	} else if exist == 0 {
		isMatch = false
	} else if Code, err2 := redisDb.Get("Register:Verification:" + email).Result(); err2 != nil {
		isMatch, err = false, err2
	} else {
		isMatch = Code == vCode
	}
	return
}

func CheckPassword(email, password string) (isMatch bool, err error) {
	if user, err1 := GetUserByEmail(email); err1 != nil {
		isMatch, err = false, err1
	} else {
		isMatch = myUtils.CheckEncryptPwdMatch(password, user.Password)
	}
	return
}
