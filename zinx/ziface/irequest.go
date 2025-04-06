package ziface

/*
 IRequest 接口定义:
	把客户端请求的链接信息,和请求的数据包封装到一个request中
*/

type IRequest interface {
	//获取链接
	GetConnection() IConnection
	//获取数据
	GetData() []byte
	//获取数据长度
	GetMsgLen() uint32
	//获取消息ID
	GetMsgID() uint32
}
