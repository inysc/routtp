package routtp

import "net/http"

func New() *Router {
	return &Router{
		NotFound: func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(404)
			w.Write([]byte("404 Not Found"))
		},
		Method: make([]Pair[string, *Node], 0, 10),
	}
}

type Router struct {
	NotFound HandlerFunc
	Method   []Pair[string, *Node]
}

func (router *Router) Add(isRoute bool, meth, path string, fn ...HandlerFunc) {
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
	root.AddRoute(isRoute, path, fn...)
}

func (router *Router) POST(path string, fn ...HandlerFunc) {
	router.Add(true, http.MethodPost, path, fn...)
}

func (router *Router) DELETE(path string, fn ...HandlerFunc) {
	router.Add(true, http.MethodDelete, path, fn...)
}

func (router *Router) PUT(path string, fn ...HandlerFunc) {
	router.Add(true, http.MethodPut, path, fn...)
}

func (router *Router) PATCH(path string, fn ...HandlerFunc) {
	router.Add(true, http.MethodPatch, path, fn...)
}

func (router *Router) GET(path string, fn ...HandlerFunc) {
	router.Add(true, http.MethodGet, path, fn...)
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
		ctx.Clean()
		ctxpool.Put(ctx)
	}()

	ctx.Request = r
	ctx.Response = w
	if !root.Get(ctx) {
		r = r.WithContext(ctx)
		router.NotFound(w, r)
	}
	r = r.WithContext(ctx)

	for ; ctx.idx < len(ctx.Fns); ctx.idx++ {
		ctx.Fns[ctx.idx](w, r)
	}
}
