package routtp

import (
	"net/http"
	"net/http/pprof"

	"github.com/inysc/facade"
)

type methods struct {
	mehtods []Pair[string, *Node]
}

func (m *methods) Append(p Pair[string, *Node]) {
	m.mehtods = append(m.mehtods, p)
}

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
		middle:  fn,
		methods: &methods{make([]Pair[string, *Node], 0, 10)},
	}
}

type Router struct {
	NotFound HandlerFunc
	middle   HandlersChain
	methods  *methods
}

func (router *Router) Add(meth, path string, fn ...HandlerFunc) {
	var root *Node
	for _, v := range router.methods.mehtods {
		if v.Key == meth {
			root = v.Val
		}
	}
	if root == nil {
		root = &Node{}
		router.methods.Append(Pair[string, *Node]{meth, root})
	}

	fns := router.combineHandlers(fn)

	facade.Debugf("[routtp] %s %s %dops", meth, path, len(fns))
	root.AddRoute(path, fns...)
}

func (router *Router) Use(fn ...HandlerFunc) {
	router.middle = append(router.middle, fn...)
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
		NotFound: router.NotFound,
		middle:   router.combineHandlers(fn),
		methods:  router.methods,
	}
}

func (router *Router) combineHandlers(fn HandlersChain) HandlersChain {
	mergedHandlers := make(HandlersChain, len(router.middle)+len(fn))
	copy(mergedHandlers, router.middle)
	copy(mergedHandlers[len(router.middle):], fn)
	return mergedHandlers
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var root *Node
	for _, v := range router.methods.mehtods {
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

func (router *Router) WithPprof() {
	router.GET("/debug/pprof", Wrap(pprof.Index))
	router.GET("/debug/pprof/:key", Wrap(pprof.Index))
	router.GET("/debug/pprof/cmdline", Wrap(pprof.Cmdline))
	router.GET("/debug/pprof/profile", Wrap(pprof.Profile))
	router.GET("/debug/pprof/symbol", Wrap(pprof.Symbol))
	router.GET("/debug/pprof/trace", Wrap(pprof.Trace))
}
