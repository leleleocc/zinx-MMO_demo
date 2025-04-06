package ziface

/*
工具包，拆包封包
 */
type IDataPack interface {
	//获取包头长度方法
	GetHeadLen() uint32
	//封包方法
	Pack(msg IMessage) ([]byte, error)
	//拆包方法
	Unpack([]byte) (IMessage, error)
}