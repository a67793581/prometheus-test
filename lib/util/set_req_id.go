package util

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SetReqId() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestId string
		requestId = c.DefaultQuery("request_id", "")
		if requestId == "" {
			requestId = c.GetHeader("X-REQID")
			if requestId == "" {
				requestId = c.GetHeader("X-TRACE-ID")
				if requestId == "" {
					requestId = uuid.New().String()
				}
			}
		}
		c.Set("request_id", requestId)
		c.Next()
	}
}

func GetRequestId(c context.Context) string {
	if c == nil {
		return "null_request_id"
	}
	tmp := c.Value("request_id")
	requestId, _ := tmp.(string)
	if requestId != "" {
		return requestId
	}
	return "null_request_id"
}
