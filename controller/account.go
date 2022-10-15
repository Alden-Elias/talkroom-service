package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon"
	"log"
	"math/rand"
	"strconv"
	"talkRoom/dao"
	"talkRoom/models"
	"talkRoom/myUtils"
)

func SentVerificationCode(c *gin.Context) {
	if email := c.Query("email"); !myUtils.IsAnEmail(email) {
		badRequest(c)
	} else if vCode, err := dao.GetVerificationCode(email); err != nil {
		serverError(c)
	} else if err = myUtils.SentVerificationCode(email, vCode); err != nil {
		serverError(c)
	} else {
		success(c, nil)
	}
}

func Register(c *gin.Context) {
	var req models.RegisterReq
	if err := c.Bind(&req); err != nil {
		badRequest(c)
	} else if !myUtils.IsAnEmail(req.Email) {
		badRequest(c)
	} else if isMatch, err := dao.VerifyCode(req.Email, req.VCode); err != nil {
		serverError(c)
	} else if !isMatch {
		forbidden(c, "验证码错误或验证码已过期！")
	} else if !myUtils.PasswordCheck(req.Password) {
		badRequest(c)
	} else if encryptPwd, err := myUtils.EncryptPwd(req.Password); err != nil {
		fmt.Println(err.Error())
		serverError(c)
	} else {
		user := models.User{
			Email:    req.Email,
			Password: encryptPwd,
		}
		if user.Avatar, err = myUtils.GetBase64AvatarByStr(strconv.Itoa(rand.Int())); err != nil {
			serverError(c)
		} else if err = dao.UserAdd(&user); err != nil {
			fmt.Println(err.Error())
			if err == dao.EmailUsed {
				forbidden(c, err.Error())
			} else {
				serverError(c)
			}
		} else if account, err := dao.GetUserByEmail(user.Email); err != nil {
			serverError(c)
		} else if tokenStr, err := myUtils.GetToken(strconv.Itoa(int(account.ID)), myUtils.UserIdentitySubject); err != nil {
			log.Println(err)
			serverError(c)
		} else {
			c.SetCookie("token", tokenStr, 0, "", "", false, true)
			success(c, struct {
				Account models.UserDetailed `json:"account"`
			}{models.UserDetailed{
				UserBrief: models.UserBrief{
					ID:     account.ID,
					Name:   account.Name,
					Avatar: account.Avatar,
					Email:  account.Email,
				},
				Description: account.Description,
				Sex:         account.Sex,
				Birthday:    carbon.Time2Carbon(account.Birthday).Format("Y-m-d"),
			}})
		}
	}
}

func Login(c *gin.Context) {
	var user models.User
	if err := c.Bind(&user); err != nil {
		badRequest(c)
	} else if isMatch, err := dao.CheckPassword(user.Email, user.Password); err != nil {
		log.Println(err)
		serverError(c)
	} else if !isMatch {
		forbidden(c, "密码错误")
	} else if account, err := dao.GetUserByEmail(user.Email); err != nil {
		log.Println(err)
		serverError(c)
	} else if tokenStr, err := myUtils.GetToken(strconv.Itoa(int(account.ID)), myUtils.UserIdentitySubject); err != nil {
		log.Println(err)
		serverError(c)
	} else {
		c.SetCookie("token", tokenStr, 0, "", "", false, true)
		success(c, struct {
			Account models.UserDetailed `json:"account"`
		}{models.UserDetailed{
			UserBrief: models.UserBrief{
				ID:     account.ID,
				Name:   account.Name,
				Avatar: account.Avatar,
				Email:  account.Email,
			},
			Description: account.Description,
			Sex:         account.Sex,
			Birthday:    carbon.Time2Carbon(account.Birthday).Format("Y-m-d"),
		}})
	}
}

func UpdateUserInfo(c *gin.Context) {
	var req struct {
		Name        string `json:"name"`
		Birthday    string `json:"birthday"`
		Sex         string `json:"sex"`
		Description string `json:"description"`
		Avatar      string `json:"avatar"`
	}

	if user, err := authorizeCheck(c); err != nil {
		unauthorized(c)
	} else if err = c.Bind(&req); err != nil {
		badRequest(c)
	} else {
		if req.Name != "" && len(req.Name) <= 35 {
			user.Name = req.Name
		}
		if req.Birthday != "" {
			user.Birthday = carbon.Parse(req.Birthday).Carbon2Time()
		}
		if req.Sex != "" && (req.Sex == "男" || req.Sex == "女") {
			user.Sex = req.Sex
		}
		if req.Description != "" && len(req.Description) <= 300 {
			user.Description = req.Description
		}
		if 0 < len(req.Avatar) && len(req.Avatar) <= 3*(1<<20) {
			user.Avatar = &req.Avatar
		}

		if err = dao.UserSave(user); err != nil {
			serverError(c)
		} else {
			var res = models.UserDetailed{
				UserBrief: models.UserBrief{
					ID:     user.ID,
					Name:   user.Name,
					Email:  user.Email,
					Avatar: user.Avatar,
				},
				Description: user.Description,
				Sex:         user.Sex,
				Birthday:    carbon.Time2Carbon(user.Birthday).Format("Y-m-d"),
			}
			success(c, struct {
				Account models.UserDetailed `json:"account"`
			}{res})
		}
	}
}
