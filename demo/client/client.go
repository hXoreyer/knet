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
