package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

/*
只负责数据的封包和解包的单元测试
*/
func TestDataPack(t *testing.T) {
	/*
		模拟的服务器
	*/
	Listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen err:", err)
		return
	}
	/*
		服务器接收数据并解包
	*/
	go func() {
		for {
			conn, err := Listener.Accept()
			if err != nil {
				fmt.Println("server accept err:", err)
				continue
			}
			go func(conn net.Conn) {
				//处理客户端请求--拆包
				dp := NewDataPack()
				for {
					headData := make([]byte, dp.GetHeadLen())
					cnt, err := io.ReadFull(conn, headData) //根据长度读conn的数据
					if cnt != 8 || err != nil {
						fmt.Println("read head error")
						break
					}
					msgHead, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println("server unpack err:", err)
						return
					}
					if msgHead.GetMsgLen() > 0 {
						//有数据，根据head中的dataLen读取data
						msg := msgHead.(*Message) //interface->struct断言
						msg.Data = make([]byte, msg.GetMsgLen())
						cnt, err := io.ReadFull(conn, msg.Data)
						if cnt != int(msg.GetMsgLen()) || err != nil {
							fmt.Println("read msg data error")
						}
						fmt.Println("==> Recv Msg: ID=", msg.GetMsgID(), ", len=", msg.GetMsgLen(), ", data=", string(msg.GetData()))
					}
				}
			}(conn)
		}
	}()

	/*
	   模拟客户端连接
	*/
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial err:", err)
		return
	}
	df := NewDataPack()
	//模拟粘包过程
	msg1 := NewMsgPackage(1, []byte{'z','i','n','x'})
	sendData1, err := df.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 err:", err)
		return
	}
	msg2 := NewMsgPackage(2, []byte{'h','e','l','l','o'})
	sendData2, err := df.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg2 err:", err)
		return
	}
	sendData1 = append(sendData1, sendData2...)
	conn.Write(sendData1)
	select {}
}
