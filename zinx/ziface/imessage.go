package ziface

/*
将请求的消息封装到Message中
*/
type IMessage interface {
	GetMsgID() uint32
	GetMsgLen() uint32
	GetData() []byte
	SetMsgID(msgID uint32)
	SetMsgLen(len uint32)
	SetData(data []byte)
}
