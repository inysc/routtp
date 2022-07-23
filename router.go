package routtp

import (
	"net/http"
)

func New(fn ...HandlerFunc) *Router {
	return &Router{
		NotFound: func(ctx *Context) {
			ctx.Exception(&exception{
				code:   http.StatusNotFound,
				status: http.StatusNotFound,
				msg:    "404 Not Found!!!",
				data:   "",
			})
		},
		Handlers: fn,
		Method:   make([]Pair[string, *Node], 0, 10),
	}
}

type Router struct {
	NotFound HandlerFunc
	Handlers HandlersChain
	Method   []Pair[string, *Node]
}

func (router *Router) Add(meth, path string, fn ...HandlerFunc) {
	var root *Node
	for _, v := range router.Method {
		if v.Key == meth {
			root = v.Val
		}
	}
	if root == nil {
		root = &Node{}
		router.Method = append(router.Method, Pair[string, *Node]{
			Key: meth,
			Val: root,
		})
	}

	fns := router.combineHandlers(fn)

	root.AddRoute(path, fns...)
}

func (router *Router) Use(fn ...HandlerFunc) {
	router.Handlers = append(router.Handlers, fn...)
}

func (router *Router) POST(path string, fn ...HandlerFunc) {
	router.Add(http.MethodPost, path, fn...)
}

func (router *Router) DELETE(path string, fn ...HandlerFunc) {
	router.Add(http.MethodDelete, path, fn...)
}

func (router *Router) PUT(path string, fn ...HandlerFunc) {
	router.Add(http.MethodPut, path, fn...)
}

func (router *Router) PATCH(path string, fn ...HandlerFunc) {
	router.Add(http.MethodPatch, path, fn...)
}

func (router *Router) GET(path string, fn ...HandlerFunc) {
	router.Add(http.MethodGet, path, fn...)
}

func (router *Router) Group(fn ...HandlerFunc) *Router {
	return &Router{
		NotFound: func(ctx *Context) {
			ctx.Exception(&exception{
				code:   http.StatusNotFound,
				status: http.StatusNotFound,
				msg:    "404 Not Found!!!",
				data:   "",
			})
		},
		Handlers: router.combineHandlers(fn),
		Method:   router.Method,
	}
}

func (router *Router) combineHandlers(fn HandlersChain) HandlersChain {
	mergedHandlers := make(HandlersChain, len(router.Handlers)+len(fn))
	copy(mergedHandlers, router.Handlers)
	copy(mergedHandlers[len(router.Handlers):], fn)
	return mergedHandlers
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var root *Node
	for _, v := range router.Method {
		if v.Key == r.Method {
			root = v.Val
		}
	}
	ctx := ctxpool.Get().(*Context)
	defer func() {
		ctx.clean()
		ctxpool.Put(ctx)
	}()

	ctx.Request = r
	ctx.Response = w
	if !root.Get(ctx, "") {
		router.NotFound(ctx)
		return
	}
	r = r.WithContext(ctx)
	ctx.Request = r
	for ; ctx.idx < len(ctx.fns); ctx.idx++ {
		ctx.fns[ctx.idx](ctx)
	}
}
