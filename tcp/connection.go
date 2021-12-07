package knet

import (
	"errors"
	"fmt"
	"io"
	"net"
)

/*
	连接模块，链接每个连接的业务
*/

type IConnection interface {
	//启动连接
	Start()
	//停止连接
	Stop()
	//获取当前连接绑定的conn
	GetConn() *net.TCPConn
	//获取当前连接的id
	GetID() uint32
	//获取远程客户端信息
	RemoteAddr() net.Addr
	//发送数据
	Send(id uint32, data []byte) error
}

type HandleFunc func(*net.TCPConn, []byte, int) error

type Connection struct {
	//连接的套接字
	Conn *net.TCPConn
	//连接的ID
	ConnID uint32
	//连接状态
	isClosed bool
	//接受停止推出状态的channel
	ExitChan chan bool
	//当前连接的处理方法
	Routers IHandler
}

//初始化模块
func NewConnection(con *net.TCPConn, id uint32, routers IHandler) IConnection {
	c := &Connection{
		Conn:     con,
		ConnID:   id,
		Routers:  routers,
		isClosed: false,
		ExitChan: make(chan bool, 1),
	}
	return c
}

//读写业务
func (c *Connection) StartReader() {
	defer c.Stop()

	for {
		//拆包
		dp := NewPack()
		head := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.Conn, head); err != nil {
			fmt.Println("read msg head err:", err)
			break
		}

		msg, err := dp.Unpack(head)
		if err != nil {
			fmt.Println("unpack err:", err)
			break
		}

		if msg.GetLen() > 0 {
			temp := make([]byte, msg.GetLen())
			if _, err := io.ReadFull(c.Conn, temp); err != nil {
				fmt.Println("read msg data err:", err)
				break
			}
			msg.SetData(temp)
		}

		//调用当前连接的路由方法
		req := &Request{
			conn: c,
			msg:  msg,
		}
		c.Routers.RunHandler(req)
	}
}

//启动连接
func (c *Connection) Start() {
	fmt.Printf("[Start] ID = %d\n", c.ConnID)
	//TODO 启动读数据业务
	go c.StartReader()
}

//停止连接
func (c *Connection) Stop() {
	fmt.Printf("[Stop] ID = %d\n", c.ConnID)

	if c.isClosed == true {
		return
	}
	c.isClosed = true
	c.Conn.Close()
	close(c.ExitChan)
}

//获取当前连接绑定的conn
func (c *Connection) GetConn() *net.TCPConn {
	return c.Conn
}

//获取当前连接的id
func (c *Connection) GetID() uint32 {
	return c.ConnID
}

//获取远程客户端信息
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//发送数据
func (c *Connection) Send(id uint32, data []byte) error {
	if c.isClosed {
		return errors.New("Connection closes")
	}

	dp := NewPack()
	buff, err := dp.Pack(NewMessage(id, data))
	if err != nil {
		fmt.Printf("pack msg err: %s ,id: %d\n", err.Error(), id)
		return err
	}

	if _, err := c.Conn.Write(buff); err != nil {
		fmt.Printf("write msg err: %s, id: %d\n", err.Error(), id)
		return err
	}
	return nil
}
