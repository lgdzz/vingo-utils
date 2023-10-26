package websocket

// JoinChanel 加入频道
func JoinChanel(channelId string, uniqueId string, isBroadcast bool) {
	// 记录用户进入的频道
	userOfChannelMutex.Lock()
	if _, ok := userOfChannel[uniqueId]; !ok {
		userOfChannel[uniqueId] = make(map[string]bool)
	}
	userOfChannel[uniqueId][channelId] = true
	userOfChannelMutex.Unlock()

	channelMutex.Lock()
	defer channelMutex.Unlock()
	// 频道内现有用户
	if channelUsers, ok := channel[channelId]; ok {
		defer func() {
			// 如果需要广播给频道其他用户
			if isBroadcast {
				for channelUserId := range channelUsers {
					Send(channelUserId, WsMessage{
						Event: "ChannelUserJoin",
						Data: map[string]string{
							"channelId": channelId,
							"userId":    uniqueId,
						},
					})
				}
			}
		}()
	} else {
		channel[channelId] = make(map[string]bool)
	}
	// 加入到频道
	channel[channelId][uniqueId] = true
}

// QuitChannel 退出频道
func QuitChannel(channelId string, uniqueId string, isBroadcast bool) {
	// 记录用户进入的频道
	userOfChannelMutex.Lock()
	defer userOfChannelMutex.Unlock()
	delete(userOfChannel[uniqueId], channelId)
	if len(userOfChannel[uniqueId]) == 0 {
		delete(userOfChannel, uniqueId)
	}

	channelMutex.Lock()
	defer channelMutex.Unlock()
	delete(channel[channelId], uniqueId)
	if len(channel[channelId]) == 0 {
		delete(channel, channelId)
	}
	if channelUsers, ok := channel[channelId]; ok {
		defer func() {
			// 如果需要广播给频道其他用户
			if isBroadcast {
				for channelUserId := range channelUsers {
					Send(channelUserId, WsMessage{
						Event: "ChannelUserQuit",
						Data: map[string]string{
							"channelId": channelId,
							"userId":    uniqueId,
						},
					})
				}
			}
		}()
	}
}

// QuitChannelAll 退出所有频道
func QuitChannelAll(uniqueId string, isBroadcast bool) {
	userOfChannelMutex.Lock()
	channelMutex.Lock()
	defer channelMutex.Unlock()
	if channels, ok := userOfChannel[uniqueId]; ok {
		userOfChannelMutex.Unlock()
		for channelId := range channels {
			delete(channel[channelId], uniqueId)
			if len(channel[channelId]) == 0 {
				delete(channel, channelId)
			}
		}
		delete(userOfChannel, uniqueId)
	}
}
