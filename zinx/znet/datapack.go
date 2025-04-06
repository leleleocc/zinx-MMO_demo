package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinx/utils"
	"zinx/ziface"
)

type DataPack struct{}

func (dp *DataPack) GetHeadLen() uint32 {
	//Message:ID uint32 4个字节 + 数据长度 uint32 4个字节()
	return 8
}

func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	dataBuffer := bytes.NewBuffer([]byte{})
	//LittleEndian 小端字节序
	if err := binary.Write(dataBuffer, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuffer, binary.LittleEndian, msg.GetMsgID()); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuffer, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return dataBuffer.Bytes(), nil
}
/*
拆包(head)
 */
func (dp *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	dataBuffer := bytes.NewReader(binaryData)
	// 读head
	msg := &Message{}
	if err := binary.Read(dataBuffer, binary.LittleEndian, &msg.MsgLen); err != nil {
		return nil, err
	}
	if err := binary.Read(dataBuffer, binary.LittleEndian, &msg.MsgID); err != nil {
		return nil, err
	}
	if msg.MsgLen > utils.GlobalObject.MaxPackageSize && utils.GlobalObject.MaxPackageSize>0 {
		return nil, errors.New("too large msg data received")
	}
	return msg, nil
}

func NewDataPack() *DataPack {
	return &DataPack{}
}
