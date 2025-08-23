package message

import (
	"github.com/gin-gonic/gin"
)

type Msg struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type MsgWithoutData struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func SendMsg(c *gin.Context, httpStatus int, msg string, data interface{}) {
	if data == nil {
		c.JSON(httpStatus, MsgWithoutData{
			Code: httpStatus,
			Msg:  msg,
		})
	} else {
		c.JSON(httpStatus, Msg{
			Code: httpStatus,
			Msg:  msg,
			Data: data,
		})
	}
}
