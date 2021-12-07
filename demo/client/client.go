package main

import (
	"fmt"
	"io"
	knet "knet/tcp"
	"net"
	"time"
)

func main() {
	con, err := net.Dial("tcp", ":5555")
	if err != nil {
		fmt.Println("dail err:", err)
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
	for {
		con.Write(buf)
		buf := make([]byte, 1024)
		n, err := con.Read(buf)

		if n == 0 {
			fmt.Printf("read err %s\n", con.RemoteAddr().String())
			break
		}
		if err != nil && err != io.EOF {
			fmt.Println("read error")
			break
		}
		fmt.Println("recv by host data = ", string(buf[:n]))

		time.Sleep(time.Second)
	}
}
