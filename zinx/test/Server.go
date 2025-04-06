package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

func (pr *PingRouter) Handle(request ziface.IRequest) {
	fmt.Printf("-->recv from client msgID: %d, data: %s\n", request.GetMsgID(), string(request.GetData()))
	if err := request.GetConnection().SendMsg(200, []byte("ping...ping...ping...")); err != nil {
		fmt.Println(err)
		return
	}
}

type PongRouter struct {
	znet.BaseRouter
}

func (pr *PongRouter) Handle(request ziface.IRequest) {
	fmt.Printf("recv from client msgID: %d, data: %s\n", request.GetMsgID(), string(request.GetData()))
	if err := request.GetConnection().SendMsg(201, []byte("pong...pong...pong...")); err != nil {
		fmt.Println(err)
		return
	}
}

func DoSomething(conn ziface.IConnection) {
	fmt.Println("=====》Do something before...")
	conn.SetProperty("name", "rifo")
	conn.SetProperty("age", 18)
}
func DoSomethingAfter(conn ziface.IConnection) {
	fmt.Println("=====》Do something after...")
	if name, err := conn.GetProperty("name"); err == nil {
		fmt.Println("name: ", name)
	}
	if age, err := conn.GetProperty("age"); err == nil {
		fmt.Println("age: ", age)
	}
}
func main() {
	s := znet.NewServer()
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &PongRouter{})

	s.SetOnConnStart(DoSomething)
	s.SetOnConnStop(DoSomethingAfter)
	s.Serve()
}
