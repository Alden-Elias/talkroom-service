package controller

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"strconv"
	"talkRoom/dao"
	"talkRoom/models"
)

const (
	systemMsgId = 0

	addFriendEventId        = 1
	cleanUnreadCountEventId = 2
)

func SentMessage(c *gin.Context) {
	var req messageReq

	if user, err := authorizeCheck(c); err != nil {
		unauthorized(c)
	} else if err = c.Bind(&req); err != nil {
		badRequest(c)
	} else if err = sentMessage(user.ID, req.To, req.Message); err != nil {
		forbidden(c, err.Error())
	} else {
		success(c, nil)
	}
}

//sentMessage 发送消息
//warning: 未检查好友关系是否成立
func sentMessage(from, to uint, message string) error {
	var (
		ws *websocket.Conn
		ok bool
	)

	if ws, ok = websockets[to]; !ok {
		return errors.New("好友未上线")
	}

	var res = messageRes{
		From:    from,
		Message: message,
	}

	//将消息存入MongoDB
	if err := dao.StoryMessage(from, to, &message); err != nil {
		return err
	}
	//warning： 未做错误处理
	if from == to {
		return nil
	}
	return ws.WriteJSON(res)
}

func GetHistory(c *gin.Context) {
	if user, err := authorizeCheck(c); err != nil {
		unauthorized(c)
	} else if to, err := strconv.Atoi(c.Query("to")); err != nil || to < 0 {
		badRequest(c)
	} else if messages, err := dao.GetMessages(user.ID, uint(to)); err != nil {
		serverError(c)
	} else {
		success(c, messages)
	}
}

func cleanUnreadCount(from, to uint) error {
	if err := dao.CleanUnreadCount(from, to); err != nil {
		return err
	} else if msg, err := json.Marshal(models.SystemMsg{
		EventId: cleanUnreadCountEventId,
		Json: struct {
			Uid uint `json:"uid"`
		}{Uid: from}}); err != nil {
		return err
	} else if err = sentMessage(systemMsgId, to, string(msg)); err != nil {
		return err
	}
	return nil
}
