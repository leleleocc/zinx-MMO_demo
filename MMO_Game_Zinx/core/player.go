package core

import (
	"MMO_Game_Zinx/pb"
	"fmt"
	"google.golang.org/protobuf/proto"
	"math/rand"
	"sync"
	"zinx/ziface"
)

type Player struct {
	Pid        int32
	Conn       ziface.IConnection //连接信息（玩家和客户端通信的connection）
	X, Y, Z, V float32            //x轴，高度，y轴，旋转角度（0-360）
}

// 生成玩家唯一id
var PIDGen int32 = 1
var IdLock sync.Mutex

func NewPlayer(conn ziface.IConnection) *Player {
	IdLock.Lock()
	defer IdLock.Unlock()
	id := PIDGen
	PIDGen++

	return &Player{
		Pid:  id,
		Conn: conn,
		X:    float32(160 + rand.Intn(10)),
		Y:    0,
		Z:    float32(150 + rand.Intn(20)),
		V:    0,
	}
}

// SendMsg /*proto序列化+发送消息*/
func (p *Player) SendMsg(msgId uint32, data proto.Message) {
	msg, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("marshal msg err:", err)
		return
	}
	if p.Conn == nil {
		fmt.Println("conn is nil")
		return
	}

	if err := p.Conn.SendMsg(msgId, msg); err != nil {
		fmt.Println("send msg err:", err)
		return
	}
	return
}

// 同步id到客户端，msgid=1
func (p *Player) SyncPid() {
	data := &pb.SyncPid{
		Pid: p.Pid,
	}
	p.SendMsg(1, data)
}

// 同步当前初始位置给客户端,msgid=200，Tp=2
func (p *Player) BroadCastStartPosition() {
	data := &pb.Broadcast{
		Pid: p.Pid,
		Tp:  2,
		Data: &pb.Broadcast_Pos{
			Pos: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}
	p.SendMsg(200, data)
}

// 世界聊天，MsgID=200,Tp=1
func (p *Player) Talk(content string) {
	data := &pb.Broadcast{
		Pid: p.Pid,
		Tp:  1,
		Data: &pb.Broadcast_Content{
			Content: content,
		},
	}
	players := WorldMgr.GetAllPlayers()
	for _, player := range players {
		player.SendMsg(200, data)
	}
}

func (p *Player) SyncPosition() {
	//得到周围玩家
	players := p.GetSurroundPlayers()
	//服务器向周围玩家发送自己的位置 MsgID=200,Tp=2
	proto_msg := &pb.Broadcast{
		Pid: p.Pid,
		Tp:  2,
		Data: &pb.Broadcast_Pos{
			Pos: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}
	for _, player := range players {
		player.SendMsg(200, proto_msg)
	}
	//服务器向自己（客户端）发送其他玩家的位置 MsgID=202
	players_proto_msg := make([]*pb.Players, 0, len(players))
	for _, player := range players {
		p := &pb.Players{
			Pid: player.Pid,
			Pos: &pb.Position{
				X: player.X,
				Y: player.Y,
				Z: player.Z,
				V: player.V,
			},
		}
		players_proto_msg = append(players_proto_msg, p)
	}
	SyncPlayers_proto_msg := &pb.SyncPlayers{
		Ps: players_proto_msg,
	}
	p.SendMsg(202, SyncPlayers_proto_msg)
}

// 广播自己的位置 MsgID=200,Tp=4
func (p *Player) UpdatePosition(x, y, z, v float32) {
	p.X, p.Y, p.Z, p.V = x, y, z, v
	proto_msg := &pb.Broadcast{
		Pid: p.Pid,
		Tp:  4,
		Data: &pb.Broadcast_Pos{
			Pos: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}
	players := p.GetSurroundPlayers()
	for _, player := range players {
		player.SendMsg(200, proto_msg)
	}
}

// 得到周围玩家
func (p *Player) GetSurroundPlayers() []*Player {
	pids := WorldMgr.AOIMgr.GetSurroundPlayersByPos(p.X, p.Z)
	players := make([]*Player, 0, len(pids))
	for _, pid := range pids {
		players = append(players, WorldMgr.GetPlayerByPid(pid))
	}
	return players
}
func (p *Player) Offline() {
	players := p.GetSurroundPlayers()
	proto_msg := &pb.SyncPid{
		Pid: p.Pid,
	}
	for _, player := range players {
		player.SendMsg(201, proto_msg)
	}
	WorldMgr.RemovePlayer(p)
	WorldMgr.AOIMgr.RemovePidFromGridByPos(p.Pid, p.X, p.Z)
}
