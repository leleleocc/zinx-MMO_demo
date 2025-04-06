package main

import (
	"MMO_Game_Zinx/apis"
	"MMO_Game_Zinx/core"
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

func OnConnStart(conn ziface.IConnection) {
	//玩家上线
	player := core.NewPlayer(conn)
	player.SyncPid()
	player.BroadCastStartPosition()
	//玩家添加到世界管理器
	core.WorldMgr.AddPlayer(player)
	//绑定pid和conn
	conn.SetProperty("pid", player.Pid)
	//同步周边玩家位置
	player.SyncPosition()
	fmt.Println("===>Player pid = ", player.Pid, "has been online.")
}
func OnConnStop(conn ziface.IConnection) {
	//玩家下线
	pid, _ := conn.GetProperty("pid")
	player := core.WorldMgr.GetPlayerByPid(pid.(int32))
	player.Offline()
	fmt.Println("===>Player pid = ", player.Pid, "has been offline.")
}
func main() {
	s := znet.NewServer()
	s.SetOnConnStart(OnConnStart)
	s.AddRouter(2, &apis.WorldChatApi{})
	s.AddRouter(3, &apis.MoveApi{})
	s.SetOnConnStop(OnConnStop)
	s.Serve()
}
