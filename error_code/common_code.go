package error_code

import (
	"github.com/gin-gonic/gin"
)

type ErrCode int32

var ClientMsgMap map[ErrCode]string

const (
	Success             ErrCode = 0
	Failed              ErrCode = -1
	Error               ErrCode = -2
	ErrorRecoveryFailed ErrCode = -3
)

var CommonClientMsg = map[ErrCode]string{
	Success:             "success",
	Failed:              "failed",
	Error:               "server_busy_please_retry",
	ErrorRecoveryFailed: "failed",
}

func GetMsg(c *gin.Context, code ErrCode) string {

	msg, ok := ClientMsgMap[code]
	if ok {
		return msg
	}
	return CommonClientMsg[Error]
}

func ErrCodeInit() {

	if ClientMsgMap == nil {
		ClientMsgMap = make(map[ErrCode]string, 0)
	}
	for k, v := range CommonClientMsg {
		ClientMsgMap[k] = v
	}

}
