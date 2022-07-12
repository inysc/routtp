package routtp

import (
	"context"
	"net/http"
	"sync"
	"time"
)

var ctxpool sync.Pool

func init() {
	ctxpool = sync.Pool{
		New: func() any {
			return &Context{
				Request:  nil,
				Response: nil,
				Param:    make([]Pair[string, string], 0, 4),
				Cancel:   func() {},
				Fns:      make(HandlersChain, 0, 4),
			}
		},
	}
}

type Pair[K, V any] struct {
	Key K
	Val V
}

type Context struct {
	Request  *http.Request
	Response http.ResponseWriter
	Param    []Pair[string, string]
	Cancel   context.CancelFunc
	Fns      HandlersChain
	idx      int
}

var _ context.Context = &Context{}

func (ctx *Context) Deadline() (deadline time.Time, ok bool) {
	return
}

func (ctx *Context) Done() <-chan struct{} {
	return nil
}

func (ctx *Context) Err() error {
	return nil
}

func (ctx *Context) Value(key any) any {
	k, ok := key.(string)
	if !ok {
		return nil
	}
	for _, v := range ctx.Param {
		if k == v.Key {
			return v.Val
		}
	}
	return nil
}

func (ctx *Context) Clone() (ctxClone *Context) {
	ctxClone = &Context{
		Request:  ctx.Request,
		Response: ctx.Response,
		Param:    make([]Pair[string, string], 0, len(ctx.Param)),
		Cancel:   func() {},
		Fns:      make(HandlersChain, 0, len(ctx.Fns)),
		idx:      ctx.idx,
	}

	copy(ctxClone.Param, ctx.Param)

	return
}

func (ctx *Context) Clean() {
	ctx.Request = nil
	ctx.Response = nil
	ctx.Param = ctx.Param[:0]
	ctx.Cancel = func() {}
	ctx.Fns = ctx.Fns[:0]
	ctx.idx = 0
}

func (ctx *Context) Prefix(path, uri string) (i, j int) {
	for ; i < len(path) && j < len(uri); i++ {
		switch path[i] {
		case ':':
			if uri[j] == '/' {
				return -1, -1
			}
			ii := i + 1
			jj := j
			for i+1 < len(path) && path[i+1] != '/' {
				i++
			}
			for j+1 < len(uri) && uri[j+1] != '/' {
				j++
			}
			ctx.Param = append(ctx.Param, Pair[string, string]{
				Key: path[ii : i+1],
				Val: uri[jj : j+1],
			})
		case '*':
			ctx.Param = append(ctx.Param, Pair[string, string]{
				Key: path[i+1:],
				Val: uri[j:],
			})
			return i + 1, j
		case uri[j]:
		default:
			return -1, -1
		}
		j++
	}
	return
}

func (ctx *Context) Next() {
	ctx.idx++
	for ctx.idx < len(ctx.Fns) {
		ctx.Fns[ctx.idx](ctx.Response, ctx.Request)
		ctx.idx++
	}
}

func (ctx *Context) Abort() {
	ctx.idx = len(ctx.Fns)
}
