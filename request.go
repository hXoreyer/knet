package knet

/*
	请求数据
	将客户的请求数据包装到Request模块中
*/
type IRequest interface {
	//获取当前连接
	GetConnection() IConnection
	//获取消息数据
	GetData() []byte
	//获取ID
	GetID() uint32
	//获取数据长度
	GetLen() uint32
	getRid() *uint32
}

type Request struct {
	//已建立的连接
	conn IConnection
	//客户端请求的数据
	msg IMessage
	//当前Request的ID
	rid *uint32
}

func NewRequest(con IConnection, msg IMessage, rid *uint32) IRequest {
	return &Request{
		conn: con,
		msg:  msg,
		rid:  rid,
	}
}

//获取当前连接
func (r *Request) GetConnection() IConnection {
	return r.conn
}

//获取消息数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

//获取ID
func (r *Request) GetID() uint32 {
	return r.msg.GetID()
}

//获取数据长度
func (r *Request) GetLen() uint32 {
	return r.msg.GetLen()
}

func (r *Request) getRid() *uint32 {
	return r.rid
}
