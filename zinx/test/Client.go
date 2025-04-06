package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx/znet"
)

func main() {
	fmt.Println("client start")
	time.Sleep(1 * time.Second)
	// 创建客户端对象
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err", err)
		return
	}
	for {
		//发
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMsgPackage(0, []byte("ping zinx server!")))
		if err != nil {
			fmt.Println("client pack err", err)
			return
		}
		if _, err := conn.Write(binaryMsg); err != nil {
			fmt.Println("client write err", err)
			return
		}
	    //收
		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("read head err", err)
			break
		}
		msgHead, err := dp.Unpack(binaryHead)
		if err != nil {
			fmt.Println("client unpack err", err)
			break
		}
		if msgHead.GetMsgLen() > 0 {
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read data err", err)
				return
			}
			fmt.Printf("-->recv server msg : id=%d, len=%d, data=%s\n", msg.GetMsgID(), msg.GetMsgLen(), msg.GetData())
		}
		//cpu阻塞
		time.Sleep(1 * time.Second)
	}
}
