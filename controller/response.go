package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"talkRoom/models"
)

//success 200成功响应
func success(c *gin.Context, data any) {
	res := models.Response{
		Status: 0,
		Data:   data,
	}
	c.JSON(http.StatusOK, res)
}

//4xx 客户端错误

//badRequest 400请求错误
func badRequest(c *gin.Context) {
	res := models.Response{
		Status: http.StatusBadRequest,
		Msg:    "请求错误无法响应",
	}
	c.JSON(http.StatusBadRequest, res)
}

//unauthorized 401无权限
func unauthorized(c *gin.Context) {
	res := models.Response{
		Status: http.StatusUnauthorized,
		Msg:    "无请求权限",
	}
	c.JSON(http.StatusUnauthorized, res)
}

//forbidden 403请求禁止
func forbidden(c *gin.Context, msg string) {
	res := models.Response{
		Status: -1,
		Msg:    msg,
	}
	c.JSON(http.StatusForbidden, res)
}

//notFound 404响应
func notFound(c *gin.Context) {
	res := models.Response{
		Status: http.StatusNotFound,
		Msg:    "未找到相关资源",
	}
	c.JSON(http.StatusNotFound, res)
}

//5xx 服务器错误

//serverError 500服务器处理错误
func serverError(c *gin.Context) {
	res := models.Response{
		Status: http.StatusInternalServerError,
		Msg:    "不可预料的错误，请联系管理员解决",
	}
	c.JSON(http.StatusInternalServerError, res)
}
