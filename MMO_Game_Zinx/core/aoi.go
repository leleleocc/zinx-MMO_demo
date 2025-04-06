package core

import "fmt"

const (
	AOI_MIN_X int = 85
	AOI_MAX_X int = 410
	AOI_MIN_Y int = 75
	AOI_MAX_Y int = 400
	AOI_CNT_X int = 3
	AOI_CNT_Y int = 3
)

type AOIManager struct {
	MinX  int           //区域左边界
	MaxX  int           //区域右边界
	MinY  int           //区域上边界
	MaxY  int           //区域下边界
	CntX  int           //x方向格子的数量
	CntY  int           //y方向格子的数量
	Grids map[int]*Grid //map存放当前区域格子信息
}

func NewAOIManager(minX, maxX, minY, maxY, cntX, cntY int) *AOIManager {
	aoiMgr := &AOIManager{
		MinX:  minX,
		MaxX:  maxX,
		MinY:  minY,
		MaxY:  maxY,
		CntX:  cntX,
		CntY:  cntY,
		Grids: make(map[int]*Grid),
	}
	//对所有格子编号
	for y := 0; y < cntY; y++ {
		for x := 0; x < cntX; x++ {
			gid := y*cntX + x
			aoiMgr.Grids[gid] = NewGrid(gid, minX+x*aoiMgr.gridWidth(), minX+(x+1)*aoiMgr.gridWidth(), minY+y*aoiMgr.gridHeight(), minY+(y+1)*aoiMgr.gridHeight())
		}
	}
	for _, grid := range aoiMgr.Grids {
		fmt.Println( grid.String())
	}
	return aoiMgr
}

// 一个格子宽度
func (aoiMgr *AOIManager) gridWidth() int {
	return (aoiMgr.MaxX - aoiMgr.MinX) / aoiMgr.CntX
}

// 一个格子高度
func (aoiMgr *AOIManager) gridHeight() int {
	return (aoiMgr.MaxY - aoiMgr.MinY) / aoiMgr.CntY
}

func (aoiMgr *AOIManager) String() string {
	s := fmt.Sprintf("AOIManager:\n MinX:%d, MaxX:%d, MinY:%d, MaxY:%d, CntX:%d, CntY:%d, Grids:\n", aoiMgr.MinX, aoiMgr.MaxX, aoiMgr.MinY, aoiMgr.MaxY, aoiMgr.CntX, aoiMgr.CntY)
	for _, grid := range aoiMgr.Grids {
		s += fmt.Sprintf("%v\n", grid)
	}
	return s
}

func (aoiMgr *AOIManager) GetSurroundGridsByGid(gid int) (grids []*Grid) {
	if _, ok := aoiMgr.Grids[gid]; !ok {
		return
	}
	grids = append(grids, aoiMgr.Grids[gid])
	//看看左右有没有格子，有就加入
	idx := gid % aoiMgr.CntX
	if idx > 0 {
		grids = append(grids, aoiMgr.Grids[gid-1])
	}
	if idx < aoiMgr.CntX-1 {
		grids = append(grids, aoiMgr.Grids[gid+1])
	}
	//看看上下有没有格子，有就加入
	for _, grid := range grids {
		idy := grid.GID / aoiMgr.CntX
		if idy > 0 {
			grids = append(grids, aoiMgr.Grids[grid.GID-aoiMgr.CntX])
		}
		if idy < aoiMgr.CntY-1 {
			grids = append(grids, aoiMgr.Grids[grid.GID+aoiMgr.CntX])
		}
	}
	return grids
}

func (aoiMgr *AOIManager) GetGidByPos(x, y float32) (gID int) {
	return (int(y)-aoiMgr.MinY)/aoiMgr.gridHeight()*aoiMgr.CntX + (int(x)-aoiMgr.MinX)/aoiMgr.gridWidth()
}

func (aoiMgr *AOIManager) GetSurroundPlayersByPos(x, y float32) (playerIDs []int32) {
	gid := aoiMgr.GetGidByPos(x, y)
	grids := aoiMgr.GetSurroundGridsByGid(gid)
	for _, grid := range grids {
		playerIDs = append(playerIDs, grid.GetPlayerIDs()...) //...表示将这个[]int打散再append
		fmt.Println("======>grid:", grid.GID, "playerIDs:", grid.GetPlayerIDs())
	}
	return playerIDs
}

func (aoiMge *AOIManager) AddPidToGrid(pID int32, gID int) {
	aoiMge.Grids[gID].AddPlayer(pID)
}
func (aoiMge *AOIManager) RemovePidFromGrid(pID int32, gID int) {
	aoiMge.Grids[gID].RemovePlayer(pID)
}
func (aoiMge *AOIManager) GetPidsByGid(gID int) (playerIDs []int32) {
	playerIDs = aoiMge.Grids[gID].GetPlayerIDs()
	return
}
func (aoiMge *AOIManager) AddPidToGridByPos(pID int32, x, y float32) {
	gID := aoiMge.GetGidByPos(x, y)
	fmt.Println("gID:",gID)
	aoiMge.AddPidToGrid(pID, gID)
}
func (aoiMge *AOIManager) RemovePidFromGridByPos(pID int32, x, y float32) {
	gID := aoiMge.GetGidByPos(x, y)
	aoiMge.RemovePidFromGrid(pID, gID)
}
