package znet

import "zinx/ziface"

// 实现router接口，先嵌入BaseRouter基类，再根据需要实现方法

type BaseRouter struct{}

func (br *BaseRouter) PreHandle(request ziface.IRequest)  {}
func (br *BaseRouter) Handle(request ziface.IRequest)     {}
func (br *BaseRouter) PostHandle(request ziface.IRequest) {}
