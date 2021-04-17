package lambda

import (
	"net/http"
	"path"
)

/**
*@Author lyer
*@Date 4/7/21 21:28
*@Describe
**/
type IRouters interface {
	Group(string, ...HandlerFunc) IRouters
	Use(...HandlerFunc) IRouters                    //加入中间件
	Handle(string, string, ...HandlerFunc) IRouters //自定义传入HTTP方法
	GET(string, ...HandlerFunc) IRouters
	POST(string, ...HandlerFunc) IRouters
	DELETE(string, ...HandlerFunc) IRouters
	HEAD(string, ...HandlerFunc) IRouters
	OPTIONS(string, ...HandlerFunc) IRouters
}

type RouterGroup struct {
	handlers HandlersChain
	prefix   string
	engine   *Engine
}

func (group *RouterGroup) handle(method string, relativePath string, handlers HandlersChain) IRouters {
	allHandlers := group.combineHandlers(handlers)
	fullPath := path.Join(group.prefix, relativePath)
	group.engine.addRoute(method, fullPath, allHandlers) //间接调用engine中root的addRoute方法
	return group
}

func (group *RouterGroup) Use(middleware ...HandlerFunc) IRouters {
	group.handlers = append(group.handlers, middleware...)
	return group
}

func (group *RouterGroup) Handle(method string, relativePath string, handlers ...HandlerFunc) IRouters {
	return group.handle(method, relativePath, handlers)
}

func (group *RouterGroup) GET(relativePath string, handlers ...HandlerFunc) IRouters {
	return group.handle(http.MethodGet, relativePath, handlers)
}

func (group *RouterGroup) POST(relativePath string, handlers ...HandlerFunc) IRouters {
	return group.handle(http.MethodPost, relativePath, handlers)
}

func (group *RouterGroup) DELETE(relativePath string, handlers ...HandlerFunc) IRouters {
	return group.handle(http.MethodDelete, relativePath, handlers)
}

func (group *RouterGroup) HEAD(relativePath string, handlers ...HandlerFunc) IRouters {
	return group.handle(http.MethodHead, relativePath, handlers)
}

func (group *RouterGroup) OPTIONS(relativePath string, handlers ...HandlerFunc) IRouters {
	return group.handle(http.MethodOptions, relativePath, handlers)
}

func (group *RouterGroup) Group(relativePath string, handlers ...HandlerFunc) IRouters {
	newPrefix := path.Join(group.prefix, relativePath)
	return &RouterGroup{
		handlers: group.combineHandlers(handlers),
		prefix:   newPrefix,
		engine:   group.engine,
	}
}

//合并handlers
func (group *RouterGroup) combineHandlers(handlers HandlersChain) HandlersChain {
	newHandlers := make(HandlersChain, len(group.handlers)+len(handlers))
	copy(newHandlers, group.handlers)
	copy(newHandlers[len(group.handlers):], handlers)
	return newHandlers
}
