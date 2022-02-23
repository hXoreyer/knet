package knet

import (
	"bytes"
	"encoding/binary"
)

/*
	封包拆包模块
*/

const headSize = 8

type IPack interface {
	//获取包头长度
	GetHeadLen() uint32
	//封包
	Pack(IMessage) ([]byte, error)
	//拆包
	Unpack([]byte) (IMessage, error)
}

type Pack struct{}

//实例化
func NewPack() IPack {
	return &Pack{}
}

//获取包头长度
func (p *Pack) GetHeadLen() uint32 {
	return headSize
}

//封包
func (p *Pack) Pack(msg IMessage) ([]byte, error) {
	//数据缓冲区
	data := &bytes.Buffer{}

	//将各个数据写进缓冲区
	if little {
		if err := binary.Write(data, binary.LittleEndian, msg.GetLen()); err != nil {
			return nil, err
		}
		if err := binary.Write(data, binary.LittleEndian, msg.GetID()); err != nil {
			return nil, err
		}
		if err := binary.Write(data, binary.LittleEndian, msg.GetData()); err != nil {
			return nil, err
		}
	} else {
		if err := binary.Write(data, binary.BigEndian, msg.GetLen()); err != nil {
			return nil, err
		}
		if err := binary.Write(data, binary.BigEndian, msg.GetID()); err != nil {
			return nil, err
		}
		if err := binary.Write(data, binary.BigEndian, msg.GetData()); err != nil {
			return nil, err
		}
	}
	return data.Bytes(), nil
}

//拆包
func (p *Pack) Unpack(data []byte) (IMessage, error) {
	//将data放入Reader里
	r := bytes.NewReader(data)

	//获取Head
	msg := &Message{}

	if little {
		if err := binary.Read(r, binary.LittleEndian, &msg.DataLen); err != nil {
			return nil, err
		}
		if err := binary.Read(r, binary.LittleEndian, &msg.Id); err != nil {
			return nil, err
		}
	} else {
		if err := binary.Read(r, binary.BigEndian, &msg.DataLen); err != nil {
			return nil, err
		}
		if err := binary.Read(r, binary.BigEndian, &msg.Id); err != nil {
			return nil, err
		}
	}
	return msg, nil
}
