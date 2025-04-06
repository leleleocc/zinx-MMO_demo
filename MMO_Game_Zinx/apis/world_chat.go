package apis

import (
	"MMO_Game_Zinx/core"
	"MMO_Game_Zinx/pb"
	"fmt"
	"google.golang.org/protobuf/proto"
	"zinx/ziface"
	"zinx/znet"
)

type WorldChatApi struct {
	znet.BaseRouter
}

func (wc *WorldChatApi) Handle(request ziface.IRequest) {
	data := &pb.Talk{}
	err := proto.Unmarshal(request.GetData(), data)
	if err != nil {
		fmt.Println("Tal unmarshal err:", err)
		return
	}
	pid, err := request.GetConnection().GetProperty("pid")
	player := core.WorldMgr.GetPlayerByPid(pid.(int32))
	//广播聊天消息
	player.Talk(data.Content)
}
