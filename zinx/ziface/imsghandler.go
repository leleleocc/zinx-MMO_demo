package ziface

/*
消息管理（msgid和路由映射）
工作池机制
*/

type IMsgHandle interface {
	//执行
	DoMsgHandler(request IRequest)
	//添加路由
	AddRouter(msgID uint32, router IRouter)
	//启动Worker工作池，只执行一次，因为一个Zinx只有一个WorkerPool
	StartWorkerPool()
	//将消息交给TaskQueue，由worker进行处理
	SendMsgToTaskQueue(request IRequest)
}
