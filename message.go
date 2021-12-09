package knet

/*
	请求消息模块
*/

type IMessage interface {
	//获取消息ID
	GetID() uint32
	//获取数据长度
	GetLen() uint32
	//获取数据
	GetData() []byte

	//设置消息ID
	SetID(uint32)
	//设置数据长度
	SetLen(uint32)
	//设置数据
	SetData([]byte)
}

type Message struct {
	Id      uint32 //消息ID
	DataLen uint32 //数据长度
	Data    []byte //数据
}

//获取消息ID
func (m *Message) GetID() uint32 {
	return m.Id
}

//获取数据长度
func (m *Message) GetLen() uint32 {
	return m.DataLen
}

//获取数据
func (m *Message) GetData() []byte {
	return m.Data
}

//设置消息ID
func (m *Message) SetID(id uint32) {
	m.Id = id
}

//设置数据长度
func (m *Message) SetLen(len uint32) {
	m.DataLen = len
}

//设置数据
func (m *Message) SetData(data []byte) {
	m.Data = data
}

//实例化

func NewMessage(id uint32, data []byte) IMessage {
	return &Message{
		Id:      id,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}
