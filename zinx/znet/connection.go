package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"zinx/utils"
	"zinx/ziface"
)

type Connection struct {
	//TCPServer
	TCPServer ziface.IServer
	//TCP套接字
	Conn *net.TCPConn
	//当前连接的ID
	ConnID uint32
	//当前连接的状态
	IsClosed bool
	//告知该连接已经退出/停止的channel（Reader告知Writer）
	ExitChan chan bool
	//该链接处理的Router方法
	MsgHandle ziface.IMsgHandle
	//读写Goroutine之间的消息通道
	MsgChan chan []byte
	//链接属性集合（用户定义）
	Property map[string]interface{}
	//链接属性保护锁
	PropertyLock sync.RWMutex
}

func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running]...ConnID=", c.ConnID)
	defer fmt.Println("[Writer is exit]	ConnID=", c.ConnID, ", remote addr is ", c.GetRemoteAddr().String())
	for {
		select {
		case data := <-c.MsgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error ", err)
				return
			}
		case <-c.ExitChan:
			return
		}
	}
}
func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running]...ConnID=", c.ConnID)
	defer fmt.Println("[Reader is exit]	ConnID=", c.ConnID, ", remote addr is ", c.GetRemoteAddr().String())
	defer c.Stop()

	for {
		//接收数据
		dp := NewDataPack()
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error ", err)
			break
		}
		msg, err := dp.Unpack(headData) //拆包
		if err != nil {
			fmt.Println("unpack error ", err)
			break
		}
		var data []byte
		if msg.GetMsgLen() > 0 {
			//有数据，则读取数据
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error ", err)
				break
			}
		}
		msg.SetData(data)
		//新建request绑定数据和conn
		req := &Request{
			conn: c,
			msg:  msg,
		}
		//处理业务
		if utils.GlobalObject.WorkerPoolSize > 0 {
			c.MsgHandle.SendMsgToTaskQueue(req)
		} else {
			//未开启工作池，直接执行
			go c.MsgHandle.DoMsgHandler(req)
		}
	}
}

func (c *Connection) Start() {
	//启动当前链接的读数据业务
	go c.StartReader()
	//启动当前链接的写数据业务
	go c.StartWriter()
	//调用Hook函数
	c.TCPServer.CallOnConnStart(c)
}
func (c *Connection) Stop() {
	fmt.Println("Connection stop...ConID=", c.ConnID)
	if c.IsClosed {
		return
	}
	c.IsClosed = true
	//调用Hook函数
	c.TCPServer.CallOnConnStop(c)
	//销毁链接
	c.Conn.Close()
	c.ExitChan <- true
	//从连接管理器中删除当前连接
	c.TCPServer.GetConnMgr().Remove(c)
	//回收资源
	close(c.MsgChan)
	close(c.ExitChan)
}
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}
func (c *Connection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}
func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	//发送数据
	if c.IsClosed {
		return errors.New("Connection closed when send msg.")
	}
	dp := NewDataPack()
	msg := NewMsgPackage(msgID, data)
	binaryMsg, err := dp.Pack(msg)
	if err != nil {
		fmt.Println("Pack error msg id = ", msgID)
		return errors.New("Pack error msg id = " + fmt.Sprintf("%d", msgID))
	}
	c.MsgChan <- binaryMsg
	return nil
}
func (c *Connection) SetProperty(key string, value interface{}) {
	c.PropertyLock.Lock()
	defer c.PropertyLock.Unlock()
	c.Property[key] = value
}
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.PropertyLock.RLock()
	defer c.PropertyLock.RUnlock()
	if value, ok := c.Property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}
func (c *Connection) RemoveProperty(key string) {
	c.PropertyLock.Lock()
	defer c.PropertyLock.Unlock()
	delete(c.Property, key)
}
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, MsgHandle ziface.IMsgHandle) *Connection {
	c := &Connection{
		TCPServer: server,
		Conn:      conn,
		ConnID:    connID,
		IsClosed:  false,
		MsgHandle: MsgHandle,
		ExitChan:  make(chan bool, 1),
		MsgChan:   make(chan []byte),
		Property:  make(map[string]interface{}),
	}
	c.TCPServer.GetConnMgr().Add(c) //将当前新创建的连接添加到ConnManager中
	fmt.Println("cur conn Num: ", c.TCPServer.GetConnMgr().Len())
	return c
}
