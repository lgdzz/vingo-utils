package websocket

import (
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/gorilla/websocket"
)

type WsMessage struct {
	Event string `json:"event"`
	Data  any    `json:"data"`
}

// Send 发送消息
func Send(uniqueId string, message WsMessage) bool {
	userConnectionsMutex.RLock()
	conn := userConnections[uniqueId]
	defer userConnectionsMutex.RUnlock()
	if conn == nil {
		return false
	}

	messageByte, err := sonic.Marshal(message)
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

// Close 断开连接
func Close(uniqueId string) {
	userConnectionsMutex.RLock()
	defer userConnectionsMutex.RUnlock()
	if conn, ok := userConnections[uniqueId]; ok {
		_ = conn.Close()
	}
}

// IsOnline 是否在线
func IsOnline(uniqueId string) bool {
	userConnectionsMutex.RLock()
	defer userConnectionsMutex.RUnlock()
	return userConnections[uniqueId] != nil
}

// OnlineUserNum 在线总用户数
func OnlineUserNum() int {
	userConnectionsMutex.RLock()
	defer userConnectionsMutex.RUnlock()
	return len(userConnections)
}

// OnlineChannelUserNum 频道内在线用户数
func OnlineChannelUserNum(channelId string) int {
	channelMutex.RLock()
	defer channelMutex.RUnlock()
	return len(channel[channelId])
}

// ChannelNum 频道数
func ChannelNum() int {
	channelMutex.RLock()
	defer channelMutex.RUnlock()
	return len(channel)
}

// UserOfChannel 用户所在频道
func UserOfChannel(uniqueId string) []string {
	userOfChannelMutex.RLock()
	defer userOfChannelMutex.RUnlock()
	if channels, ok := userOfChannel[uniqueId]; ok {
		result := make([]string, 0)
		for channelId := range channels {
			result = append(result, channelId)
		}
		return result
	}
	return []string{}
}
