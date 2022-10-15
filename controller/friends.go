package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"strconv"
	"talkRoom/dao"
	"talkRoom/models"
	"talkRoom/myUtils"
)

//FriendAdd 添加好友
func FriendAdd(c *gin.Context) {
	if user, err := authorizeCheck(c); err != nil {
		log.Println(err)
		unauthorized(c)
	} else if uid, err := strconv.Atoi(c.Query("uid")); err != nil {
		badRequest(c)
	} else if isFriend := dao.IsFriend(user.ID, uint(uid)); isFriend {
		forbidden(c, "你们已经是好友了，快开始聊天吧！")
	} else if s, err := json.Marshal(models.FromAndTo{
		From: user.ID,
		To:   uint(uid),
	}); err != nil {
		log.Println(err)
		serverError(c)
	} else if token, err := myUtils.GetToken(string(s), myUtils.AddFriendsSubject); err != nil {
		serverError(c)
	} else if msg, err := json.Marshal(models.SystemMsg{EventId: addFriendEventId, Json: struct {
		Uname string `json:"uname"`
		Uid   uint   `json:"uid"`
		Token string `json:"token"`
	}{Uname: user.Name, Uid: user.ID, Token: token}}); err != nil {
		serverError(c)
	} else if err = sentMessage(systemMsgId, uint(uid), string(msg)); err != nil {
		forbidden(c, err.Error())
	} else {
		success(c, nil)
	}
}

func ConfirmFriend(c *gin.Context) {
	var fromAndTo models.FromAndTo
	if user, err := authorizeCheck(c); err != nil {
		unauthorized(c)
	} else if token := c.Query("token"); token == "" {
		badRequest(c)
	} else if _, claims, err := myUtils.ParseToken(token); err != nil {
		log.Println(err)
		badRequest(c)
	} else if claims.Subject != myUtils.AddFriendsSubject {
		badRequest(c)
	} else if err := json.Unmarshal([]byte(claims.Identity), &fromAndTo); err != nil {
		forbidden(c, "token解析错误")
	} else if user.ID != fromAndTo.To {
		unauthorized(c)
	} else if err := dao.AddFriend(fromAndTo.From, fromAndTo.To); err != nil {
		log.Println(err)
		serverError(c)
	} else {
		success(c, nil)
	}
}

func ListFriends(c *gin.Context) {
	//var friends []models.UserBrief
	if user, err := authorizeCheck(c); err != nil {
		unauthorized(c)
	} else if friends, err := dao.ListFriendsBrief(user.ID); err != nil {
		if err == gorm.ErrRecordNotFound {
			success(c, nil)
		} else {
			serverError(c)
		}
	} else {
		success(c, struct {
			Users *[]models.UserBrief `json:"users"`
		}{friends})
	}
}
