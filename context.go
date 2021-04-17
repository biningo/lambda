package lambda

import "net/http"

/**
*@Author lyer
*@Date 3/25/21 13:49
*@Describe
**/
type Context struct {
	cwriter  responseWriter      //自己实现的接口 扩展了一些方法
	Writer   http.ResponseWriter //原始的writer
	Request  *http.Request
	Path     string
	engine   *Engine
	Method   string
	handlers HandlersChain
	Params   []Param
	index    int
}

func (c *Context) Init(w http.ResponseWriter, req *http.Request) {
	c.Request = req
	c.Writer = w
	c.cwriter.ResponseWriter = w
	c.Path = c.Request.URL.Path
	c.Method = c.Request.Method
	c.index = -1
}

func (c *Context) Next() {
	c.index++
	for c.index < len(c.handlers) {
		c.handlers[c.index](c)
		c.index++
	}
}

func (c *Context) Status(code int) {
	c.Writer.WriteHeader(code)
}

type responseWriter struct {
	http.ResponseWriter
}
