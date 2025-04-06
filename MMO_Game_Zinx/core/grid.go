package core

import (
	"fmt"
	"sync"
)

type Grid struct {
	GID       int          //格子ID
	MinX      int          //格子左边界
	MaxX      int          //格子右边界
	MinY      int          //格子下边界
	MaxY      int          //格子上边界
	playerIDs map[int32]bool //格子内玩家集合
	pidLock   sync.RWMutex
}

func NewGrid(GID, MinX, MaxX, MinY, MaxY int) *Grid {
	return &Grid{
		GID:       GID,
		MinX:      MinX,
		MaxX:      MaxX,
		MinY:      MinY,
		MaxY:      MaxY,
		playerIDs: make(map[int32]bool),
	}
}

func (g *Grid) AddPlayer(pID int32) {
	g.pidLock.Lock()
	defer g.pidLock.Unlock()
	g.playerIDs[pID] = true
}
func (g *Grid) RemovePlayer(pID int32) {
	g.pidLock.Lock()
	defer g.pidLock.Unlock()
	delete(g.playerIDs, pID)
}

func (g *Grid) GetPlayerIDs() (playerIDs []int32) {
	g.pidLock.RLock()
	defer g.pidLock.RUnlock()
	for pID := range g.playerIDs {
		playerIDs = append(playerIDs, pID)
	}
	return
}

func (g *Grid) String() string {
	return fmt.Sprintf("Grid id: %d, minX: %d, maxX: %d, minY: %d, maxY: %d, playerIDs: %v", g.GID, g.MinX, g.MaxX, g.MinY, g.MaxY, g.playerIDs)
}

