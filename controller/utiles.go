package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
	"talkRoom/dao"
	"talkRoom/models"
	"talkRoom/myUtils"
)

var (
	subjectFault      = errors.New("token主题错误")
	identityTypeFault = errors.New("identity类型错误")
)

func authorizeCheck(c *gin.Context) (user *models.User, err error) {
	if token, err1 := c.Cookie("token"); err1 != nil {
		err = err1
	} else if _, claims, err1 := myUtils.ParseToken(token); err1 != nil {
		err = err1
	} else if claims.Subject != myUtils.UserIdentitySubject {
		err = subjectFault
	} else if id, err1 := strconv.Atoi(claims.Identity); err1 != nil || id < 0 {
		err = identityTypeFault
	} else {
		user, err = dao.GetUserById(uint(id))
	}
	return
}
