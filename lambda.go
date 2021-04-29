package lambda

import (
	"net/http"
	"sync"
)

/**
*@Author lyer
*@Date 3/25/21 13:36
*@Describe
**/
type HandlerFunc func(*Context)
type HandlersChain []HandlerFunc
type methodTrees []methodTree
type H map[string]interface{}
type Engine struct {
	RouterGroup
	trees       methodTrees //每种HTTP请求方法都会有一颗trie树
	contextPool sync.Pool   //复用Context对象  此对象每个HTTP请求都会生成一个
}

func (trees methodTrees) get(method string) *node {
	for _, v := range trees {
		if v.method == method {
			return v.root
		}
	}
	return &node{}
}

func New() *Engine {
	engine := &Engine{
		trees: make([]methodTree, 0, 9), //预先分配cap=9 有9个方法
	}
	engine.RouterGroup = RouterGroup{
		prefix: "/",
		engine: engine,
	}
	engine.contextPool.New = func() interface{} {
		return &Context{engine: engine}
	}
	return engine
}

func (engine *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	context := engine.contextPool.Get().(*Context)

	//Context Init
	context.Init(w, r)
	context.engine = engine
	//handler
	engine.httpRequestHandler(context)
	engine.contextPool.Put(context)
}

func (engine *Engine) httpRequestHandler(c *Context) {
	root := engine.trees.get(c.Method)
	c.handlers, c.Params = root.getHandlers(c.Path)
	if len(c.handlers) == 0 {
		c.Writer.WriteHeader(404)
		c.Writer.Write([]byte("404"))
		return
	}
	c.Next()
}

//调用trie添加路由
func (engine *Engine) addRoute(method string, path string, handlers HandlersChain) {
	root := engine.trees.get(method)
	engine.trees = append(engine.trees, methodTree{method: method, root: root})
	root.addRoute(path, handlers)
}
