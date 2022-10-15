package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"talkRoom/models"
)

var websockets map[uint]*websocket.Conn

func init() {
	websockets = make(map[uint]*websocket.Conn)
}

const (
	readMessageEventId = 1
)

type messageRes struct {
	From    uint   `json:"from"`
	Message string `json:"message"`
}

type messageReq struct {
	To      uint   `json:"to"`
	EventId uint   `json:"eventId"`
	Message string `json:"message"`
}

//WsConnect websocket连接处理
//warning: 当消息传输失败时是不会有任何反馈的，需要在后续改进版本
//warning: 未经测试
func WsConnect(c *gin.Context) {
	var (
		user *models.User
		err  error
	)
	if user, err = authorizeCheck(c); err != nil {
		unauthorized(c)
		return
	}
	//服务升级，对于来到的http连接进行服务升级，升级到ws
	upGrande := websocket.Upgrader{
		//设置允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		//设置请求协议
		Subprotocols: []string{c.GetHeader("Sec-WebSocket-Protocol")},
	}
	//创建连接
	cn, err := upGrande.Upgrade(c.Writer, c.Request, nil)
	defer cn.Close()
	if err != nil {
		panic(err)
	}

	//一个用户只能建立一个ws连接
	if ws, ok := websockets[user.ID]; ok {
		ws.Close()
	}

	//连接存储到map中
	websockets[user.ID] = cn
	defer delete(websockets, user.ID)

	//设置连接关闭处理器，当连接主动关闭时需要退出死循环
	connected := true
	cn.SetCloseHandler(func(code int, text string) error {
		connected = false
		return nil
	})

	//接收消息并处理
	for connected {
		var req messageReq
		if err = cn.ReadJSON(&req); err != nil {
			log.Println(err)
			continue
		}

		switch req.EventId {
		case readMessageEventId: //对方消息以读，清空未读消息计数
			if err := cleanUnreadCount(user.ID, req.To); err != nil {
				log.Println(user.ID, ": readMessageEvent has an error ->", err)
			}
		}
	}
}
