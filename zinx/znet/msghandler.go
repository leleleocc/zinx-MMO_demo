package znet

import (
	"fmt"
	"zinx/utils"
	"zinx/ziface"
)

type MsgHandle struct {
	// 存放MsgId和对应的handle
	Apis         map[uint32]ziface.IRouter
	TaskQueue    []chan ziface.IRequest
	WorkPoolSize uint32
}

func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID = ", request.GetMsgID(), "is not found！Register first！")
		return
	}
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}
func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	if _, exists := mh.Apis[msgID]; exists {
		panic("repeated api msgID = " + fmt.Sprintf("%d", msgID))
	}
	mh.Apis[msgID] = router
	fmt.Println("-->Add api MsgId = ", msgID, "successfully")
}

func (mh *MsgHandle) StartWorkerPool() {
	for i := 0; i < int(mh.WorkPoolSize); i++ {
		//启动worker
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen) //指定chan最大容量
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID = ", workerID, "is started")
	for {
		select {
		//有消息过来，一个客户端的request出列，执行request绑定的业务
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	//根据客户端的connID来分配当前连接到哪个worker中（简单的负载均衡）
	workerID := request.GetConnection().GetConnID() % mh.WorkPoolSize
	fmt.Println("-->Add ConnID = ", request.GetConnection().GetConnID(),
		"request msgID = ", request.GetMsgID(),
		"to workerID = ", workerID)
	mh.TaskQueue[workerID] <- request
}
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:         make(map[uint32]ziface.IRouter),
		TaskQueue:    make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
		WorkPoolSize: utils.GlobalObject.WorkerPoolSize,
	}
}
