package knet

import "net"

type hookHandler func(c IConnection)
type overloadHandler func(c *net.TCPConn)

type IServer interface {
	//运行服务器
	Run()
	//启动服务器
	start()
	//停止服务器
	Stop()

	//添加路由
	AddRouter(id uint32, router IRouter)
	//路由三部曲
	Before(uint32, RouterFunc)
	On(uint32, RouterFunc)
	After(uint32, RouterFunc)
	//添加全局中间件
	Use(RouterFunc)
	Abort()

	//设置工作池大小
	SetWorkPoolSize(uint32)
	GetRid() *uint32

	//获取连接管理
	GetManager() IManager

	//超出连接数回调函数
	OverLoad(overloadHandler)
	SetMaxCon(size uint32)

	//连接创建hook
	OnStart(hookHandler)
	runOnStart(IConnection)
	//连接断开hook
	OnStop(hookHandler)
	runOnStop(IConnection)
}
