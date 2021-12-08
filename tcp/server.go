package knet

import (
	"fmt"
	"net"
)

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
}

type TCPServer struct {
	IP        string
	Port      int
	Conn      *net.TCPListener
	IPVersion string
	Routers   IHandler
	Version   string
	Name      string
	rids      *uint32
}

//启动服务器
func (s *TCPServer) start() {
	//获取TCP地址
	addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("[Error] Resolve tcp addr err:", err)
		return
	}

	//监听服务器地址
	s.Conn, err = net.ListenTCP(s.IPVersion, addr)
	if err != nil {
		fmt.Println("[Error] Listen tcp err:", err)
		return
	}

	//监听成功输出
	fmt.Println("[Init] Start server success...")
	//开启工作池
	s.Routers.RunWorkPool()
	//循环接受用户连接
	for {
		con, err := s.Conn.AcceptTCP()
		if err != nil {
			fmt.Println("[Error] Accept err:", err)
			continue
		}
		var cid uint32
		cid = 0
		dealCon := NewConnection(con, cid, s.Routers, s.rids)
		cid++
		go dealCon.Start()
	}

}

//运行服务器
func (s *TCPServer) Run() {
	s.start()
}

//停止服务器
func (s *TCPServer) Stop() {
	s.Conn.Close()
}

//添加路由
func (s *TCPServer) AddRouter(id uint32, router IRouter) {
	s.Routers.AddRouter(id, router)
}

//添加全局中间件
func (s *TCPServer) Use(rf RouterFunc) {
	s.Routers.Use(rf)
}

//终断请求
func (s *TCPServer) Abort() {
	s.Routers.Abort()
}

//创建新的Server模块
func NewTCPServer(ip string, port int) IServer {
	return &TCPServer{
		IPVersion: "tcp4",
		IP:        ip,
		Port:      port,
		Routers:   NewHandler(),
		Name:      "Knet",
		Version:   "V1.0",
		rids:      new(uint32),
	}
}

/*
	非接口方法
*/

//设置名称
func (s *TCPServer) SetName(name string) {
	s.Name = name
}

//设置版本号
func (s *TCPServer) SetVersion(ver string) {
	s.Version = ver
}

// Hook
func (s *TCPServer) Before(id uint32, rf RouterFunc) {
	s.Routers.Before(id, rf)
}
func (s *TCPServer) On(id uint32, rf RouterFunc) {
	s.Routers.On(id, rf)
}
func (s *TCPServer) After(id uint32, rf RouterFunc) {
	s.Routers.After(id, rf)
}

//设置工作池大小
func (s *TCPServer) SetWorkPoolSize(size uint32) {
	s.Routers.SetWorkPoolSize(size)
}
