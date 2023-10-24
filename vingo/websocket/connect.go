package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/lgdzz/vingo-utils/vingo"
	"net/http"
	"sync"
)

var (
	// 记录所用用户的websocket连接
	userConnections      = make(map[string]*websocket.Conn)
	userConnectionsMutex = sync.RWMutex{}
	// 记录所有用户的频道ID
	userOfChannel      = make(map[string]map[string]bool)
	userOfChannelMutex = sync.RWMutex{}
	// 记录所有频道，以及频道中的用户ID
	channel      = make(map[string]map[string]bool)
	channelMutex = sync.RWMutex{}
)

// 默认值
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 设置upgrader值
func WsSetUpgrader(u websocket.Upgrader) {
	upgrader = u
}

// 客户端请求连接并绑定唯一用户ID
func WsConnect(c *vingo.Context, uniqueId string, handle func(message string, uniqueId string)) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		panic(err.Error())
	}

	defer func() {
		// 连接关闭时，从映射中移除用户
		userConnectionsMutex.RLock()
		delete(userConnections, uniqueId)
		userConnectionsMutex.RUnlock()
		QuitChannelAll(uniqueId, true)
		_ = conn.Close()
	}()

	userConnectionsMutex.RLock()
	userConnections[uniqueId] = conn
	userConnectionsMutex.RUnlock()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			return
		}
		switch messageType {
		case websocket.TextMessage:
			// 处理文本消息
			if handle != nil {
				handle(string(p), uniqueId)
			}
		case websocket.BinaryMessage:
			// 处理二进制消息
		case websocket.CloseMessage:
			// 处理关闭消息
		case websocket.PingMessage:
			// 处理 ping 消息
		case websocket.PongMessage:
			// 处理 pong 消息
		default:
			// 未知消息类型
		}
		//fmt.Println(messageType, string(p))
		// Handle WebSocket messages here
		// You can send messages using conn.WriteMessage()
	}
}
