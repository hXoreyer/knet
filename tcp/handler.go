package knet

import (
	"fmt"
)

/*
	消息管理模块
*/

type IHandler interface {
	//调度Router
	RunHandler(request IRequest)
	//添加处理逻辑
	AddRouter(id uint32, router IRouter)

	//hook三部曲
	Before(uint32, RouterFunc)
	On(uint32, RouterFunc)
	After(uint32, RouterFunc)
	//添加全局中间件
	Use(RouterFunc)
}

type Handler struct {
	routers     map[uint32]IRouter
	Middlewares []RouterFunc
}

//调度Router
func (h *Handler) RunHandler(request IRequest) {
	handle, ok := h.routers[request.GetID()]
	if !ok {
		fmt.Println("use unadded handler, id = ", request.GetID())
		return
	}
	for k := range h.Middlewares {
		h.Middlewares[k](request)
	}
	handle.start(request)
}

//添加处理逻辑
func (h *Handler) AddRouter(id uint32, router IRouter) {
	h.routers[id] = router
	fmt.Println("[Handler] Add ID = ", id)
}

//hook三部曲
func (h *Handler) Before(id uint32, rf RouterFunc) {
	_, ok := h.routers[id]
	if !ok {
		h.routers[id] = NewRouter()
	}
	h.routers[id].Before(rf)
	fmt.Printf("[Handler] Add ID = %d, type = %s\n", id, "Brefore")
}
func (h *Handler) On(id uint32, rf RouterFunc) {
	_, ok := h.routers[id]
	if !ok {
		h.routers[id] = NewRouter()
	}
	h.routers[id].On(rf)
	fmt.Printf("[Handler] Add ID = %d, type = %s\n", id, "On")
}
func (h *Handler) After(id uint32, rf RouterFunc) {
	_, ok := h.routers[id]
	if !ok {
		h.routers[id] = NewRouter()
	}
	h.routers[id].After(rf)
	fmt.Printf("[Handler] Add ID = %d, type = %s\n", id, "After")
}
func (h *Handler) Use(rf RouterFunc) {
	h.Middlewares = append(h.Middlewares, rf)
}

//实例化
func NewHandler() IHandler {
	return &Handler{
		routers:     make(map[uint32]IRouter),
		Middlewares: make([]RouterFunc, 0),
	}
}
