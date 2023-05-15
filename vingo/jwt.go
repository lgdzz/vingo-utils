package vingo

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type JwtTicket struct {
	Key string `json:"key"`
	TK  string `json:"tk"`
}

type JwtBody struct {
	ID       string     `json:"id"`
	Day      uint       `json:"day"` // 默认有效期90天
	Business any        `json:"business"`
	CheckTK  bool       `json:"checkTk"`
	Ticket   *JwtTicket `json:"ticket"`
}

// 生成token
func JwtIssued(body JwtBody, signingKey string) string {
	if body.Day == 0 {
		body.Day = 90
	}
	day := 3600 * 24 * int64(body.Day)
	exp := time.Now().Unix() + day
	if body.CheckTK {
		body.Ticket = &JwtTicket{Key: MD5(fmt.Sprintf("%v%v", signingKey, body.ID)), TK: RandomString(50)}
		Redis.Set(body.Ticket.Key, body.Ticket.TK, time.Second*time.Duration(day))
	}
	signedString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"id": body.ID, "checkTk": body.CheckTK, "ticket": body.Ticket, "business": body.Business, "exp": exp}).SignedString([]byte(signingKey))
	if err != nil {
		panic(err)
	}
	return signedString
}

// 验证token
func JwtCheck(token string, signingKey string) JwtBody {
	claims, err := jwt.ParseWithClaims(token, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(signingKey), nil
	})
	if err != nil {
		panic(&AuthException{Message: err.Error()})
	}
	var body JwtBody
	CustomOutput(claims.Claims, &body)
	if body.CheckTK {
		if body.Ticket.TK != RedisResult(Redis.Get(body.Ticket.Key)) {
			panic(&AuthException{Message: "登录已失效"})
		}
	}
	return body
}
