package apis

import (
	"MMO_Game_Zinx/core"
	"MMO_Game_Zinx/pb"
	"fmt"
	"google.golang.org/protobuf/proto"
	"zinx/ziface"
	"zinx/znet"
)

type MoveApi struct {
	znet.BaseRouter
}

func (m *MoveApi) Handle(request ziface.IRequest) {
	proto_msg := &pb.Position{}
	fmt.Println("moving!!")
	err := proto.Unmarshal(request.GetData(), proto_msg)
	if err != nil {
		fmt.Println("Move position unmarshal error")
		return
	}
	pid, err := request.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("GetProperty pid error")
		return
	}
	player := core.WorldMgr.GetPlayerByPid(pid.(int32))
	player.UpdatePosition(proto_msg.X, proto_msg.Y, proto_msg.Z, proto_msg.V)
}


