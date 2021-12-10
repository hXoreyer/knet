# Knet TCP框架

基于Golang的轻量级并发TCP框架

#### [中文文档](https://hxoreyer.github.io)

##### 实例:

server:

```go
package main

import (
	"fmt"
	"net"
	"time"

	"github.com/hxoreyer/knet"
)

func main() {
	logger := knet.NewKlog("./logger")
	logger.SetUpdateTime("01:00:00") //日志文件更换时间设置，默认为"00:00:00"

	//日志三种格式，分别在success.log,info.log,error.log文件里
	logger.Success("Create logger success!!!")
	logger.Info("INFO INFO INFO!!!")
	logger.Error("Create logger error!!!")

	s := knet.NewTCPServer("127.0.0.1", 5555)
	//设置最大连接数
	s.SetMaxCon(1)

	//设置工作池数量
	s.SetWorkPoolSize(10)

	//全局中间件
	s.Use(func(request knet.IRequest) {
		fmt.Println("[Middleware] This is middleware1, Id:", request.GetID())
	})
	s.Use(func(request knet.IRequest) {
		fmt.Println("[Middleware] This is middleware2, Id:", request.GetID())
		if string(request.GetData()) == "hxoreyer" {
			fmt.Println("[Middleware] Abort from middleware2")
			s.Abort()
		}
	})

	//请求路由
	s.Before(1, func(request knet.IRequest) {
		fmt.Printf("[Router] Recv Before, ID = %d\n", request.GetID())
	})

	s.On(1, func(request knet.IRequest) {
		fmt.Printf("[Router] Recv from %s, ID = %d Data = %s\n", request.GetConnection().RemoteAddr().String(), request.GetID(), request.GetData())
		request.GetConnection().Send(request.GetID(), request.GetData())
	})
	s.On(2, func(request knet.IRequest) {
		fmt.Printf("[Router] Recv from %s, ID = %d Data = %s\n", request.GetConnection().RemoteAddr().String(), request.GetID(), request.GetData())
		request.GetConnection().Send(request.GetID(), request.GetData())
	})

	s.After(2, func(request knet.IRequest) {
		fmt.Printf("[Router] Recv After, ID = %d\n", request.GetID())
	})

	//超出最大连接数
	s.OverLoad(func(c *net.TCPConn) {
		dp := knet.NewPack()

		msg1 := &knet.Message{
			Id:      9,
			DataLen: 8,
			Data:    []byte("overload"),
		}
		buf, _ := dp.Pack(msg1)
		c.Write(buf)
		time.Sleep(time.Millisecond)
	})
	//连接开始时
	s.OnStart(func(c knet.IConnection) {
		c.SetProperty("name", "keing")
	})
	//连接结束时
	s.OnStop(func(c knet.IConnection) {
		fmt.Println(c.GetProperty("name"))
	})
	s.Run()
}

```

client:

```go
package main

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/hxoreyer/knet"
)

var exit = false

func main() {
	con, err := net.Dial("tcp", ":5555")
	if err != nil {
		fmt.Println("dail err:", err)
		return
	}
	go Reader(con)
	dp := knet.NewPack()

	msg1 := &knet.Message{
		Id:      1,
		DataLen: 5,
		Data:    []byte("keing"),
	}
	buf, _ := dp.Pack(msg1)

	msg2 := &knet.Message{
		Id:      2,
		DataLen: 8,
		Data:    []byte("hxoreyer"),
	}
	buf2, _ := dp.Pack(msg2)
	buf = append(buf, buf2...)
	for {
		con.Write(buf)
		time.Sleep(time.Second)
		if exit {
			break
		}
	}
}

func Reader(con net.Conn) {
	for {
		dp := knet.NewPack()
		head := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(con, head); err != nil {
			fmt.Println("read msg head err:", err)
			exit = true
			break
		}

		msg, err := dp.Unpack(head)
		if err != nil {
			fmt.Println("unpack err:", err)
			break
		}

		if msg.GetLen() > 0 {
			temp := make([]byte, msg.GetLen())
			if _, err := io.ReadFull(con, temp); err != nil {
				fmt.Println("read msg data err:", err)
				break
			}
			msg.SetData(temp)
		}
		fmt.Printf("[Recv] ID = %d, Data = %s\n", msg.GetID(), string(msg.GetData()))
	}
}

```

