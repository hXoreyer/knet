package knet

/*
	路由模块
	数据都是IRequest请求
*/
type IRouter interface {
	//业务前主后三种方法
	Before(router RouterFunc)
	On(router RouterFunc)
	After(router RouterFunc)
	//开启业务
	start(request IRequest)
}

type Router struct {
	middlewares []RouterFunc
	before      RouterFunc
	on          RouterFunc
	after       RouterFunc
}

//实例化
func NewRouter() IRouter {
	return &Router{
		before: func(request IRequest) {},
		on:     func(request IRequest) {},
		after:  func(request IRequest) {},
	}
}

func (r *Router) Before(router RouterFunc) {
	r.before = router
}
func (r *Router) On(router RouterFunc) {
	r.on = router
}
func (r *Router) After(router RouterFunc) {
	r.after = router
}
func (r *Router) Use(router RouterFunc) {
	r.middlewares = append(r.middlewares, router)
}
func (r *Router) start(request IRequest) {
	r.before(request)
	r.on(request)
	r.after(request)
}

type RouterFunc func(request IRequest)
