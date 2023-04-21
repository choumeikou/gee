package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request
	// request info
	Path   string
	Method string
	// response info
	StatusCode int
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}

// 从POST的body中读取key的value
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// 从URL中读取key的value
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// 设置响应的http状态码
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// 设置四种类型的响应内容，均采用字节切片方式传输
// 响应的正确调用顺序是Header().Set,然后是WriteHeader()，最后是Write()
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil { //这里存在问题，如果err!=nil，http.Error不会生效，因为前面已经按顺序完成了Header().Set、WriteHeader()和Write()，gin框架render/json.go#L56中依靠return来规避这一系列操作
		http.Error(c.Writer, err.Error(), 500)
	}
}

// 这个没有SetHeader是因为默认格式就是二进制数据传输
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}
