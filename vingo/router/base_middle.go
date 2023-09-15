package router

import (
	"github.com/gin-gonic/gin"
	"time"
)

func BaseMiddle(hook *Hook) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Set("requestStart", startTime)

		if hook.BaseMiddle != nil {
			hook.BaseMiddle(c)
		}
	}
}
