package ziface

/*
 连接管理模块抽象接口
*/

type IConnManager interface {
	//添加连接
	Add(conn IConnection)
	//删除连接
	Remove(conn IConnection)
	//根据ConnID获取连接
	Get(connID uint32) (IConnection, error)
	//得到当前连接总数
	Len() int
	//清除所有连接
	ClearConn()
}
