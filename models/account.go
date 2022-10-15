package models

import (
	"github.com/golang-jwt/jwt"
)

//RegisterReq 注册账户请求
type RegisterReq struct {
	Email    string `json:"email"`    // 邮箱，用户注册用邮箱
	Password string `json:"password"` // 密码，用户密码
	VCode    string `json:"vCode"`    // 验证码，用户验证码
}

//Claims 声明， 保存在token中的字段
type Claims struct {
	Identity string
	jwt.StandardClaims
}
