package znet

type Message struct {
	MsgLen uint32
	MsgID  uint32
	Data   []byte
}

func (m *Message) GetMsgID() uint32 {
	return m.MsgID
}
func (m *Message) GetMsgLen() uint32 {
	return m.MsgLen
}
func (m *Message) GetData() []byte {
	return m.Data
}
func (m *Message) SetMsgID(msgID uint32) {
	m.MsgID = msgID
}
func (m *Message) SetMsgLen(MsgLen uint32) {
	m.MsgLen = MsgLen
}
func (m *Message) SetData(data []byte) {
	m.Data = data
}

func NewMsgPackage(msgID uint32, data []byte) *Message {
	return &Message{
		MsgID: msgID,
		Data:  data,
		MsgLen: uint32(len(data)),
	}
}
