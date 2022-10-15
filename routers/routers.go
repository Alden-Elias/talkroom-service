package routers

import (
	"github.com/gin-gonic/gin"
	"talkRoom/controller"
)

func Start() {
	r := gin.Default()
	//r.Static("/static", "./static")
	//r.StaticFS("/imgs", http.Dir("./imgs"))
	//r.LoadHTMLGlob("templates/pages/*")

	//r.GET("/", func(c *gin.Context) {
	//	c.HTML(http.StatusOK, "chat.html", nil)
	//})
	//r.GET("/login", func(c *gin.Context) {
	//	c.HTML(http.StatusOK, "login.html", nil)
	//})
	//r.GET("/register", func(c *gin.Context) {
	//	c.HTML(http.StatusOK, "register.html", nil)
	//})
	//r.GET("/test", func(c *gin.Context) {
	//	c.HTML(http.StatusOK, "test.html", nil)
	//})

	//用户相关
	user := r.Group("/user")
	user.POST("/add", controller.UserAdd)
	user.GET("/search", controller.SearchUsers)

	//账号相关
	account := r.Group("/account")
	account.GET("/verificationCode", controller.SentVerificationCode)
	account.POST("/register", controller.Register)
	account.POST("/login", controller.Login)
	account.PUT("/updateUserInfo", controller.UpdateUserInfo)

	//好友相关
	friends := r.Group("/friends")
	friends.POST("/add", controller.FriendAdd)
	friends.POST("/confirm", controller.ConfirmFriend)
	friends.GET("/list", controller.ListFriends)

	//websocket相关
	websocket := r.Group("/websocket")
	websocket.GET("/message", controller.WsConnect)

	//massage相关
	message := r.Group("message")
	message.POST("/sent", controller.SentMessage)
	message.GET("/getHistory", controller.GetHistory)

	if r.Run() != nil {
		panic("Gin engin failed run")
	}
}
