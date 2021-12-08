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
	//终断请求
	Abort()

	//工作池
	//启动工作池
	RunWorkPool()
	//向工作池发送请求
	Send2Tasks(rq IRequest)
	//设置工作池大小
	SetWorkPoolSize(size uint32)
}

type Handler struct {
	routers      map[uint32]IRouter
	Middlewares  []RouterFunc
	abort        bool
	workpoolSize uint32
	tasks        []chan IRequest
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
		if h.abort {
			h.abort = false
			return
		}
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

//全局中间件
func (h *Handler) Use(rf RouterFunc) {
	h.Middlewares = append(h.Middlewares, rf)
}

//终断请求
func (h *Handler) Abort() {
	h.abort = true
}

//工作池数量
func (h *Handler) SetWorkPoolSize(size uint32) {
	h.workpoolSize = size
}

//实例化
func NewHandler() IHandler {
	return &Handler{
		routers:      make(map[uint32]IRouter),
		Middlewares:  make([]RouterFunc, 0),
		abort:        false,
		workpoolSize: 10,
	}
}

//启动工作池
func (h *Handler) RunWorkPool() {
	h.tasks = make([]chan IRequest, h.workpoolSize)
	for i := 0; i < int(h.workpoolSize); i++ {
		h.tasks[i] = make(chan IRequest)
		go h.runWork(h.tasks[i])
	}
	fmt.Printf("[WorkPool] %d workpools are Running...\n\n", h.workpoolSize)
}

//启动一个工作
func (h *Handler) runWork(tr chan IRequest) {
	for {
		select {
		case rq := <-tr:
			h.RunHandler(rq)
			rid := rq.getRid()
			(*rid)--
		}
	}
}

func (h *Handler) Send2Tasks(rq IRequest) {
	id := *(rq.getRid()) % h.workpoolSize
	fmt.Printf("[WorkPool] Task id:%d work in pool id:%d\n", rq.GetID(), id)
	h.tasks[id] <- rq
}
