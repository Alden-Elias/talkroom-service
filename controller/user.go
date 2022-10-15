package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"talkRoom/dao"
	"talkRoom/models"
	"talkRoom/myUtils"
)

func UserAdd(c *gin.Context) {
	var user models.User
	if err := c.Bind(&user); err != nil {
		badRequest(c)
	} else if !myUtils.IsAnEmail(user.Email) || !myUtils.PasswordCheck(user.Password) {
		badRequest(c)
	} else if encryptPwd, err := myUtils.EncryptPwd(user.Password); err != nil {
		fmt.Println(err.Error())
		serverError(c)
	} else {
		user.Password = encryptPwd
		if err = dao.UserAdd(&user); err != nil {
			fmt.Println(err.Error())
			if err == dao.EmailUsed {
				forbidden(c, err.Error())
			} else {
				serverError(c)
			}
		} else {
			success(c, nil)
		}
	}
}

func SearchUsers(c *gin.Context) {
	if _, err := authorizeCheck(c); err != nil {
		unauthorized(c)
	} else if query := c.Query("query"); /* query == "" */ false { //暂时不对请求进行限制  允许用户搜索所有用户
		badRequest(c)
	} else if users, err := dao.UserFuzzySearch(query); err != nil {
		notFound(c)
	} else {
		success(c, struct {
			Users []models.UserBrief `json:"users"`
		}{Users: *users})
	}
}
