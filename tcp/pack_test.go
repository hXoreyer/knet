package knet

import (
	"fmt"
	"io"
	"net"
	"sync"
	"testing"
)

func TestPack(t *testing.T) {
	ls, err := net.Listen("tcp", ":5520")
	if err != nil {
		fmt.Println("listen err:", err)
		return
	}
	var wg sync.WaitGroup
	go func() {
		for {
			con, err := ls.Accept()
			if err != nil {
				fmt.Println("server accept err:", err)
			}
			wg.Add(1)
			go func(con net.Conn) {
				defer wg.Done()
				dp := NewPack()

				for {
					head := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(con, head)
					if err != nil {
						break
					}
					msgHead, _ := dp.Unpack(head)
					if msgHead.GetLen() > 0 {
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetLen())

						io.ReadFull(con, msg.Data)

						fmt.Printf("Recv ID:%d, Len:%d, Data: %s\n", msg.GetID(), msg.GetLen(), string(msg.GetData()))
					}
				}
			}(con)
		}
	}()

	c, err := net.Dial("tcp", "127.0.0.1:5520")
	if err != nil {
		fmt.Println("dial err:", err)
		return
	}
	dp := NewPack()

	msg1 := &Message{
		Id:      1,
		DataLen: 5,
		Data:    []byte("keing"),
	}
	buf, _ := dp.Pack(msg1)

	msg2 := &Message{
		Id:      2,
		DataLen: 8,
		Data:    []byte("hxoreyer"),
	}
	buf2, _ := dp.Pack(msg2)
	buf = append(buf, buf2...)
	c.Write(buf)
	wg.Wait()
}
