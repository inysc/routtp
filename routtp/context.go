package routtp

import (
	"context"
	"net/http"
	"time"
)

type Param struct {
	Key string
	Val string
}

type Context struct {
	Request  *http.Request
	Response http.ResponseWriter
	Param    []Param
	Cancel   context.CancelFunc
	Fns      HandlersChain
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
		Param:    make([]Param, 0, len(ctx.Param)),
		Cancel:   func() {},
		Fns:      make(HandlersChain, 0, len(ctx.Fns)),
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
}
