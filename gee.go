package gee

import (
	"net/http"
)

// type HandlerFunc func(http.ResponseWriter, *http.Request) day1的方式，还没封装Context
type HandlerFunc func(*Context)

type Engine struct {
	// router map[string]HandlerFunc day-1，未封装
	router *router
}

func New() *Engine {
	// return &Engine{router: make(map[string]HandlerFunc)}
	return &Engine{router: newRouter()}
}

func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	// key := method + "-" + pattern
	// engine.router[key] = handler day-1

	engine.router.addRoute(method, pattern, handler)
}

func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// key := req.Method + "-" + req.URL.Path
	// if handler, ok := engine.router[key]; ok {
	// 	handler(w, req)
	// } else {
	// 	fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	// }

	c := newContext(w, req)
	engine.router.handle(c)
}
