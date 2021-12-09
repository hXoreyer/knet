package knet

import (
	"fmt"
	"knet"
	"net"
	"sync"
	"testing"
)

func TestPack(t *testing.T) {
	var wg sync.WaitGroup
	go func() {
		s := knet.NewTCPServer("127.0.0.1", 5520)
		s.On(1, func(request knet.IRequest) {
			fmt.Printf("[Router] Recv from %s, ID = %d Data = %s\n", request.GetConnection().RemoteAddr().String(), request.GetID(), request.GetData())
			request.GetConnection().Send(request.GetID(), request.GetData())
		})
		s.On(2, func(request knet.IRequest) {
			fmt.Printf("[Router] Recv from %s, ID = %d Data = %s\n", request.GetConnection().RemoteAddr().String(), request.GetID(), request.GetData())
			request.GetConnection().Send(request.GetID(), request.GetData())
		})
		s.Run()
	}()

	c, err := net.Dial("tcp", "127.0.0.1:5520")
	if err != nil {
		fmt.Println("dial err:", err)
		return
	}
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
	c.Write(buf)
	wg.Wait()
}
