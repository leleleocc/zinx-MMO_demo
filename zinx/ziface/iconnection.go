package ziface

/*
 连接模块抽象层
*/

import "net"

type IConnection interface {
	// 启动连接，接收数据
	Start()
	//停止连接，结束当前连接状态M
	Stop()
	//获取当前连接的绑定socket conn
	GetTCPConnection() *net.TCPConn
	//获取当前连接模块的连接ID
	GetConnID() uint32
	//获取远程客户端的TCP状态 IP port
	GetRemoteAddr() net.Addr
	//发送数据
	SendMsg(msgId uint32, data []byte) error
	//设置链接属性
	SetProperty(key string, value interface{})
	//获取链接属性
	GetProperty(key string) (interface{}, error)
	//移除链接属性
	RemoveProperty(key string)
}
