package models

type Response struct {
	Msg    string `json:"msg"`    // 错误提示，错误提示
	Status int    `json:"status"` // 错误码，非0表示错误
	Data   any    `json:"data"`   //响应主体
}

type FromAndTo struct {
	From uint `json:"from"` //用户请求的发送方
	To   uint `json:"to"`   //用户请求的接受方
}
