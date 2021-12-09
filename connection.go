package knet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
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
	//连接属性
	SetProperty(key string, val interface{})
	GetProperty(key string) (interface{}, error)
	DelProperty(key string)
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
	//无缓冲的读写channel
	dataChan chan []byte
	//父server
	father IServer
	//保护属性
	property map[string]interface{}
	//保护锁
	propertyLock sync.RWMutex
}

//初始化模块
func NewConnection(f IServer, con *net.TCPConn, id uint32, routers IHandler) IConnection {
	c := &Connection{
		Conn:     con,
		ConnID:   id,
		Routers:  routers,
		isClosed: false,
		ExitChan: make(chan bool, 1),
		dataChan: make(chan []byte),
		father:   f,
		property: make(map[string]interface{}),
	}
	c.father.GetManager().Add(c)
	return c
}

//写业务
func (c *Connection) StartWriter() {
	for {
		select {
		case data := <-c.dataChan:
			//读写channel有数据时
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("[Error] Send data err:", err)
				return
			}
		case <-c.ExitChan:
			return
		}
	}
}

//读业务
func (c *Connection) StartReader() {
	defer c.Stop()

	for {
		//拆包
		dp := NewPack()
		head := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.Conn, head); err != nil {
			fmt.Println("[Error] Read msg head err:", err)
			break
		}

		msg, err := dp.Unpack(head)
		if err != nil {
			fmt.Println("[Error] unpack err:", err)
			break
		}

		if msg.GetLen() > 0 {
			temp := make([]byte, msg.GetLen())
			if _, err := io.ReadFull(c.Conn, temp); err != nil {
				fmt.Println("[Error] read msg data err:", err)
				break
			}
			msg.SetData(temp)
		}

		//调用当前连接的路由方法
		rid := c.father.GetRid()
		req := NewRequest(c, msg, rid)
		*(rid)++
		c.Routers.Send2Tasks(req)
	}
}

//启动连接
func (c *Connection) Start() {
	fmt.Printf("[In] ID = %d, Addr = %s\n", c.ConnID, c.GetConn().RemoteAddr().String())
	//启动读数据业务
	go c.StartReader()
	// 启动写数据业务
	go c.StartWriter()

	go c.father.runOnStart(c)
}

//停止连接
func (c *Connection) Stop() {
	fmt.Printf("[Out] ID = %d, Addr = %s\n", c.ConnID, c.GetConn().RemoteAddr().String())

	if c.isClosed == true {
		return
	}
	c.isClosed = true
	c.father.runOnStop(c)
	c.Conn.Close()
	c.ExitChan <- true
	close(c.ExitChan)
	close(c.dataChan)
	c.father.GetManager().Del(c)
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
		fmt.Printf("[Error] pack msg err: %s ,id: %d\n", err.Error(), id)
		return err
	}

	c.dataChan <- buff
	return nil
}

//连接属性
func (c *Connection) SetProperty(key string, val interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = val
}
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if val, ok := c.property[key]; ok {
		return val, nil
	}
	return nil, errors.New("No Property")
}
func (c *Connection) DelProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}
