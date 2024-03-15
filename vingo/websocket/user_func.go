package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
)

type WsMessage struct {
	Event string `json:"event"`
	Data  any    `json:"data"`
}

// Send 发送消息
func Send(uniqueId string, message WsMessage) bool {
	userConnectionsMutex.Lock()
	conn := userConnections[uniqueId]
	defer userConnectionsMutex.Unlock()
	if conn == nil {
		return false
	}

	messageByte, err := json.Marshal(message)
	if err != nil {
		panic(err.Error())
		return false
	}

	err = conn.WriteMessage(websocket.TextMessage, messageByte)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

// IsOnline 是否在线
func IsOnline(uniqueId string) bool {
	userConnectionsMutex.Lock()
	defer userConnectionsMutex.Unlock()
	return userConnections[uniqueId] != nil
}

// OnlineUserNum 在线总用户数
func OnlineUserNum() int {
	userConnectionsMutex.Lock()
	defer userConnectionsMutex.Unlock()
	return len(userConnections)
}

// OnlineChannelUserNum 频道内在线用户数
func OnlineChannelUserNum(channelId string) int {
	channelMutex.Lock()
	defer channelMutex.Unlock()
	return len(channel[channelId])
}

// ChannelNum 频道数
func ChannelNum() int {
	channelMutex.Lock()
	defer channelMutex.Unlock()
	return len(channel)
}

// UserOfChannel 用户所在频道
func UserOfChannel(uniqueId string) []string {
	userOfChannelMutex.Lock()
	defer userOfChannelMutex.Unlock()
	if channels, ok := userOfChannel[uniqueId]; ok {
		result := make([]string, 0)
		for channelId := range channels {
			result = append(result, channelId)
		}
		return result
	}
	return []string{}
}
