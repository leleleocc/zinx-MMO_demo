package ziface

/*
IServer接口定义
*/

type IServer interface {
	// 启动服务器
	Start()
	// 停止服务器
	Stop()
	// 运行服务器
	Serve()
	//路由功能
	AddRouter(MsgId uint32, router IRouter)
	//得到当前Server的连接管理器
	GetConnMgr() IConnManager
	//设置该Server的连接创建时Hook函数
	SetOnConnStart(func(connection IConnection))
	//设置该Server的连接断开时的Hook函数
	SetOnConnStop(func(connection IConnection))
	//调用连接OnConnStart Hook函数
	CallOnConnStart(conn IConnection)
	//调用连接OnConnStop Hook函数
	CallOnConnStop(conn IConnection)
}
