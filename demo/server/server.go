package main

import (
	"fmt"
	knet "knet/tcp"
)

func main() {
	s := knet.NewTCPServer("127.0.0.1", 5555)
	s.SetWorkPoolSize(2)
	/*
		//全局中间件
		s.Use(func(request knet.IRequest) {
			fmt.Println("[Middleware] This is middleware1 by Id:", request.GetID())
		})
		s.Use(func(request knet.IRequest) {
			fmt.Println("[Middleware] This is middleware2 by Id:", request.GetID())
			if string(request.GetData()) == "hxoreyer" {
				fmt.Println("[Middleware] Abort for middleware2")
				s.Abort()
			}
		})

		//请求路由
		s.Before(1, func(request knet.IRequest) {
			fmt.Printf("[Router] Recv Before, ID = %d\n", request.GetID())
		})
	*/
	s.On(1, func(request knet.IRequest) {
		fmt.Printf("[Router] Recv from %s, ID = %d Data = %s\n", request.GetConnection().RemoteAddr().String(), request.GetID(), request.GetData())
		request.GetConnection().Send(request.GetID(), request.GetData())
	})
	s.On(2, func(request knet.IRequest) {
		fmt.Printf("[Router] Recv from %s, ID = %d Data = %s\n", request.GetConnection().RemoteAddr().String(), request.GetID(), request.GetData())
		request.GetConnection().Send(request.GetID(), request.GetData())
	})
	/*
		s.After(2, func(request knet.IRequest) {
			fmt.Printf("[Router] Recv After, ID = %d\n", request.GetID())
		})
	*/
	s.Run()
}
