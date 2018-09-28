package logic

import "server/msg"

var (
	playerInfos = make(map[uint]*PlayerInfo)
)

func init() {
	player := playerInfos[0]
	player.WriteMsg(&msg.LoginFaild{Code: msg.LoginFaild_InnerError})
	player.Close()
}
