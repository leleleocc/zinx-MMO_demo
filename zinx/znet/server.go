package znet

import (
	"fmt"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int
	//工作池和handler映射
	MsgHandle ziface.IMsgHandle
	//conn管理器
	ConnMgr ziface.IConnManager
	//hook函数
	OnConnStart func(conn ziface.IConnection)
	OnConnStop  func(conn ziface.IConnection)
}

func (s *Server) Start() {
	fmt.Println("+----------------+---------------+")
	fmt.Println("|      参数      |       值       |")
	fmt.Println("+----------------+---------------+")
	fmt.Printf("| %-14s | %-13s |\n", "Name", s.Name)
	fmt.Printf("| %-14s | %-13s |\n", "IPVersion", s.IPVersion)
	fmt.Printf("| %-14s | %-13s |\n", "Host", utils.GlobalObject.Host)
	fmt.Printf("| %-14s | %-13d |\n", "Port", s.Port)
	fmt.Printf("| %-14s | %-13d |\n", "MaxConn", utils.GlobalObject.MaxConn)
	fmt.Printf("| %-14s | %-13d |\n", "MaxPackageSize", utils.GlobalObject.MaxPackageSize)
	fmt.Printf("| %-14s | %-13d |\n", "WorkerPoolSize", utils.GlobalObject.WorkerPoolSize)
	fmt.Printf("| %-14s | %-13d |\n", "MaxWorkerTaskLen", utils.GlobalObject.MaxWorkerTaskLen)
	fmt.Println("+----------------+---------------+")

	go func() {
		//0、启动worker工作池
		s.MsgHandle.StartWorkerPool()
		//1、获取Tcp的一个addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error:", err)
			return
		}
		//2、监听服务
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, "err", err)
			return
		}
		fmt.Println("start Zinx server success:", s.Name, "listening on", s.IP, s.Port)
		var cid uint32
		cid = 0
		//3、阻塞等待客户端连接，处理客户端连接业务
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				fmt.Println("Too Many Connection MaxConn = ", utils.GlobalObject.MaxConn)
				conn.Close()
				continue
			}
			// 将链接的业务与conn进行绑定，得到链接模块
			dealConn := NewConnection(s, conn, cid, s.MsgHandle)
			cid++
			//CallBackApi在curConn的start()里调用
			go dealConn.Start()
		}
	}()
}
func (s *Server) Stop() {
	//TODO Server.Stop()
	fmt.Println("[Server Stop]")
	s.ConnMgr.ClearConn()
}
func (s *Server) Serve() {
	s.Start() ///start()函数使用go func，防止等待客户端连接阻塞start(),进而阻塞主进程，而主进程一把会同时进行其他任务如心跳检测等
	//TODO Server.Serve()
	select {} //start()函数使用go func，若不阻塞，则主进程执行完毕就退出，导致start()函数中的协程无法执行
}
func (s *Server) AddRouter(MsgId uint32, router ziface.IRouter) {
	s.MsgHandle.AddRouter(MsgId, router)
}
func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}
func (s *Server) SetOnConnStart(hookFunc func(conn ziface.IConnection)) {
	s.OnConnStart = hookFunc
}
func (s *Server) SetOnConnStop(hookFunc func(conn ziface.IConnection)) {
	s.OnConnStop = hookFunc
}
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		s.OnConnStart(conn)
	}
}
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		s.OnConnStop(conn)
	}
}
func NewServer() ziface.IServer {
	s := &Server{
		Name:      utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		MsgHandle: NewMsgHandle(),
		ConnMgr:   NewConnManager(),
	}
	return s
}
