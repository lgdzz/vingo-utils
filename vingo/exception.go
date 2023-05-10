package vingo

import (
	"github.com/gin-gonic/gin"
)

type AuthException struct {
	Message string
}

// DbException 数据库事务异常
type DbException struct {
	Message string
}

// ExceptionHandler 异常处理
func ExceptionHandler(c *gin.Context) {
	context := &Context{c}

	defer func() {
		if err := recover(); err != nil {
			switch t := err.(type) {
			case string:
				context.Response(&ResponseData{Message: t, Status: 200, Error: 1, ErrorType: "业务错误"})
			case *DbException:
				context.Response(&ResponseData{Message: t.Message, Status: 200, Error: 1, ErrorType: "数据库错误"})
			case *AuthException:
				context.Response(&ResponseData{Message: t.Message, Status: 401, Error: 1})
			default:
				context.Response(&ResponseData{Message: t.(error).Error(), Status: 200, Error: 1, ErrorType: "异常错误"})
			}
			c.Abort()
		}
	}()
	c.Next()
}
