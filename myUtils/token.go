package myUtils

import (
	"github.com/golang-jwt/jwt"
	"github.com/golang-module/carbon"
	"talkRoom/models"
	"talkRoom/setting"
)

const (
	UserIdentitySubject = "用户身份认证"
	AddFriendsSubject   = "添加用户请求"
)

var (
	jwtKey = []byte(setting.Config.Webset.KwtKey)
)

//GetToken 通过主题和身份验证获取token
func GetToken(identity string, subject string) (tokenStr string, err error) {
	now := carbon.Now()
	expiresAt := now.AddHours(6)

	claims := models.Claims{
		Identity: identity,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt.Carbon2Time().Unix(),
			IssuedAt:  now.Carbon2Time().Unix(),
			Issuer:    "root",
			Subject:   subject,
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err = t.SignedString(jwtKey)
	return
}

func ParseToken(tokenStr string) (*jwt.Token, *models.Claims, error) {
	claim := models.Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, &claim, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, nil
	})
	return token, &claim, err
}
