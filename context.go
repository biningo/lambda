package lambda

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
)

/**
*@Author lyer
*@Date 3/25/21 13:49
*@Describe
**/
type Context struct {
	cwriter    responseWriter      //自己实现的接口 扩展了一些方法
	Writer     http.ResponseWriter //原始的writer
	Request    *http.Request
	Path       string
	StatusCode int
	engine     *Engine
	Method     string
	handlers   HandlersChain
	Params     map[string]string
	index      int
	cache      map[string]interface{}
	mu         sync.RWMutex
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

func (c *Context) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	val, exists := c.cache[key]
	c.mu.Unlock()
	return val, exists
}

func (c *Context) Set(key string, val interface{}) {
	c.mu.Lock()
	if c.cache == nil {
		c.cache = make(map[string]interface{})
	}
	c.cache[key] = val
	c.mu.Unlock()
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}
func (c *Context) GetHeader(key string) string {
	return c.Request.Header.Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}

func (c *Context) BindJSON(v interface{}) error {
	if c.Request.Header.Get("Content-Type") != "application/json" {
		return errors.New("content type error")
	}
	buf := make([]byte, c.Request.ContentLength)
	c.Request.Body.Read(buf)
	return json.Unmarshal(buf, v)
}

func (c *Context) PostForm(key string) string {
	return c.Request.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

func (c *Context) Param(key string) string {
	return c.Params[key]
}

type responseWriter struct {
	http.ResponseWriter
}
