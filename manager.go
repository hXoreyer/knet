package knet

import (
	"errors"
	"sync"
)

/*
	连接管理模块
*/

type IManager interface {
	//添加连接
	Add(con IConnection)
	//删除连接
	Del(con IConnection)
	//获取连接
	Get(id uint32) (IConnection, error)
	//获取连接数量
	Num() int
	//关闭所有连接
	Clear()
}

type Manager struct {
	connections map[uint32]IConnection
	lock        sync.RWMutex
}

//实例化
func NewManager() IManager {
	return &Manager{
		connections: make(map[uint32]IConnection),
	}
}

//添加连接
func (m *Manager) Add(con IConnection) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.connections[con.GetID()] = con
}

//删除连接
func (m *Manager) Del(con IConnection) {
	m.lock.Lock()
	defer m.lock.Unlock()

	delete(m.connections, con.GetID())
}

//获取连接
func (m *Manager) Get(id uint32) (IConnection, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if con, ok := m.connections[id]; ok {
		return con, nil
	}
	return nil, errors.New("connection not found")
}

//获取连接数量
func (m *Manager) Num() int {
	return len(m.connections)
}

//关闭所有连接
func (m *Manager) Clear() {
	m.lock.Lock()
	defer m.lock.Unlock()

	for id, con := range m.connections {
		con.Stop()
		delete(m.connections, id)
	}
}
