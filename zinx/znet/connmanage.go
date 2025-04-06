package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx/ziface"
)

type ConnManager struct {
	//管理的连接map
	connections map[uint32]ziface.IConnection
	//保护连接map的读写锁
	connLock sync.RWMutex
}

func (cm *ConnManager) Add(conn ziface.IConnection) {
	//写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	cm.connections[conn.GetConnID()] = conn
	fmt.Println("connection add to ConnManager successfully: connID = ", conn.GetConnID())
}
func (cm *ConnManager) Remove(conn ziface.IConnection) {
	//写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	delete(cm.connections, conn.GetConnID())
	fmt.Println("connection Remove from ConnManager successfully: connID = ", conn.GetConnID())
}
func (cm *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	//读锁
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()
	if conn, exists := cm.connections[connID]; !exists {
		return nil, errors.New("connection not exists")
	} else {
		return conn, nil
	}
}
func (cm *ConnManager) Len() int {
	return len(cm.connections)
}
func (cm *ConnManager) ClearConn() {
	// todo 存在死锁问题
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	for connID, conn := range cm.connections {
		conn.Stop()
		delete(cm.connections, connID)
		fmt.Println("Clear all connections")
	}
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}
