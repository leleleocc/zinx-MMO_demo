package utils

import (
	"encoding/json"
	"os"
	"zinx/ziface"
)

/*
存储一切有关Zinx框架的全局参数，供其他模块使用
*/
type GlobalObj struct {
	/*
	   Server
	*/
	TcpServer ziface.IServer
	Host      string
	TcpPort   int
	Name      string

	/*
	   Zinx
	*/
	Version        string
	MaxConn        int
	MaxPackageSize uint32 //TCP数据包最大值
	WorkerPoolSize uint32 //当前业务工作worker池的数量
	MaxWorkerTaskLen uint32 //每个队列中最大的任务数量
}

var GlobalObject *GlobalObj

func (g *GlobalObj) Reload() {
	data, err := os.ReadFile("config/zinx.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}
func init() {
	GlobalObject = &GlobalObj{
		Name:           "Zinx v0.4",
		Version:        "V0.4",
		Host:           "0.0.0.0",
		TcpPort:        8999,
		MaxConn:        1000,
		MaxPackageSize: 4096,
		WorkerPoolSize: 10,
		MaxWorkerTaskLen: 1024,
	}
	//从config/zinx.json文件中加载一些配置
	GlobalObject.Reload()
}
