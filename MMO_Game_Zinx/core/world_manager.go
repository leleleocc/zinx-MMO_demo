package core

import "sync"

type WorldManager struct {
	AOIMgr  *AOIManager
	Players map[int32]*Player
	pLock   sync.RWMutex
}

// 初始化全局唯一的世界管理器
var WorldMgr *WorldManager

func init() {
	WorldMgr = &WorldManager{
		//创建地图
		AOIMgr: NewAOIManager(AOI_MIN_X, AOI_MAX_X, AOI_MIN_Y, AOI_MAX_Y, AOI_CNT_X, AOI_CNT_Y),
		//初始化玩家
		Players: make(map[int32]*Player),
	}
}
func (WorldMgr *WorldManager) AddPlayer(player *Player) {
	WorldMgr.pLock.Lock()
	WorldMgr.Players[player.Pid] = player
	WorldMgr.pLock.Unlock()

	//将玩家添加到格子中
	WorldMgr.AOIMgr.AddPidToGridByPos(player.Pid, player.X, player.Z)
}

func (WorldMgr *WorldManager) RemovePlayer(player *Player) {
	WorldMgr.pLock.Lock()
	delete(WorldMgr.Players, player.Pid)
	WorldMgr.pLock.Unlock()

	//将玩家从格子中移除
	WorldMgr.AOIMgr.RemovePidFromGridByPos(player.Pid, player.X, player.Y)
}

func (WorldMgr *WorldManager) GetPlayerByPid(pid int32) *Player {
	WorldMgr.pLock.RLock()
	defer WorldMgr.pLock.RUnlock()
	if player, ok := WorldMgr.Players[pid]; ok {
		return player
	} else {
		return nil
	}
}

func (WorldMgr *WorldManager) GetAllPlayers() []*Player {
	WorldMgr.pLock.RLock()
	defer WorldMgr.pLock.RUnlock()
	players := make([]*Player, 0)
	for _, player := range WorldMgr.Players {
		players = append(players, player)
	}
	return players
}
