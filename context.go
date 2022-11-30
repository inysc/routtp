package routtp

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/inysc/facade"
)

var ctxpool sync.Pool

func init() {
	ctxpool = sync.Pool{
		New: func() any {
			return &Context{
				Request:  nil,
				Response: nil,
				param:    make([]Pair[string, string], 0, 4),
				Cancel:   func() {},
				fns:      make(HandlersChain, 0, 4),
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
	param    []Pair[string, string]
	values   []Pair[string, any]
	Cancel   context.CancelFunc
	fns      HandlersChain
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

func (ctx *Context) Param(key string) string {
	for _, v := range ctx.param {
		if key == v.Key {
			return v.Val
		}
	}
	return ""
}

// not implate
func (ctx *Context) Value(key any) any {
	for _, v := range ctx.values {
		if v.Key == key {
			return v.Val
		}
	}
	return nil
}

func (ctx *Context) Clone() (ctxClone *Context) {
	ctxClone = &Context{
		Request:  ctx.Request,
		Response: ctx.Response,
		param:    make([]Pair[string, string], 0, len(ctx.param)),
		Cancel:   func() {},
		fns:      make(HandlersChain, 0, len(ctx.fns)),
		idx:      ctx.idx,
	}

	copy(ctxClone.param, ctx.param)

	return
}

func (ctx *Context) clean() {
	ctx.Request = nil
	ctx.Response = nil
	ctx.param = ctx.param[:0]
	ctx.Cancel = func() {}
	ctx.fns = ctx.fns[:0]
	ctx.idx = 0
}

func (ctx *Context) prefix(path, uri string) (i, j int) {
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
			ctx.param = append(ctx.param, Pair[string, string]{
				Key: path[ii : i+1],
				Val: uri[jj : j+1],
			})
		case '*':
			ctx.param = append(ctx.param, Pair[string, string]{
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
	for ctx.idx < len(ctx.fns) {
		ctx.fns[ctx.idx](ctx)
		ctx.idx++
	}
}

func (ctx *Context) Abort() {
	ctx.idx = len(ctx.fns)
}

// ---------- Request ----------
// ---------- Request Header ----------
func (ctx *Context) HeaderGet(key string) string { return ctx.Request.Header.Get(key) }

func (ctx *Context) HeaderSet(key, value string) { ctx.Request.Header.Set(key, value) }

func (ctx *Context) HeaderAdd(key, value string) { ctx.Request.Header.Add(key, value) }

func (ctx *Context) HeaderDel(key string) { ctx.Request.Header.Del(key) }

func (ctx *Context) HeaderValues(key string) []string { return ctx.Request.Header.Values(key) }

func (ctx *Context) HeaderWrite(w io.Writer) error { return ctx.Request.Header.Write(w) }

func (ctx *Context) HeaderWriteSubset(w io.Writer, exclude map[string]bool) error {
	return ctx.Request.Header.WriteSubset(w, exclude)
}

func (ctx *Context) HeaderClone() http.Header { return ctx.Request.Header.Clone() }

// ---------- Request Body ----------
func (ctx *Context) GetBody() []byte {
	buff := bytes.NewBuffer(make([]byte, 0, ctx.Request.ContentLength))
	io.Copy(buff, ctx.Request.Body)
	ctx.Request.Body.Close()
	return buff.Bytes()
}

func (ctx *Context) RealIP() string {
	ip := ctx.HeaderGet("x-apigw-ip")
	if ip != "" {
		return ip
	}

	ip = ctx.HeaderGet("X-Real-IP")
	if ip != "" {
		return ip
	}

	return ctx.Request.RemoteAddr
}

func (ctx *Context) UserPrior() uint8 {
	up := ctx.HeaderGet("x-apigw-user-prior")
	t, err := strconv.ParseUint(up, 10, 8)
	if err != nil {
		facade.Errorf("parse user prior failed, err: %v", err)
	}
	return uint8(t)
}

func (ctx *Context) Username() string { return ctx.HeaderGet("x-apigw-username") }

func (ctx *Context) Traceid() string { return ctx.HeaderGet("x-apigw-traceid") }

func (ctx *Context) Userid() string { return ctx.HeaderGet("x-apigw-userid") }

func (ctx *Context) City() string { return ctx.HeaderGet("x-apigw-city") }

func (ctx *Context) BindJSON(v any) error {
	return jsonbinding.Bind(ctx.Request, v)
}

func (ctx *Context) BindQuery(v any) error {
	return mapForm(v, ctx.Request.Form)
}

func (ctx *Context) Bind(v any) error {
	if strings.Contains(ctx.HeaderGet("Content-Type"), "application/json") {
		return ctx.BindJSON(v)
	}
	return ctx.BindQuery(v)
}

// ---------- Response ----------
func (ctx *Context) Write(p []byte) (int, error) { return ctx.Response.Write(p) }

func (ctx *Context) Header() http.Header { return ctx.Response.Header() }

func (ctx *Context) WriteHeader(statusCode int) { ctx.Response.WriteHeader(statusCode) }

func (ctx *Context) Set(key string, val any) {
	ctx.values = append(ctx.values, Pair[string, any]{key, val})
}

func (ctx *Context) Get(key string) any {
	for _, v := range ctx.values {
		if v.Key == key {
			return v.Val
		}
	}
	return nil
}

func (ctx *Context) GetBool(key string) bool {
	for _, v := range ctx.values {
		if v.Key == key {
			val, ok := v.Val.(bool)
			if ok {
				return val
			}
		}
	}
	return false
}

func (ctx *Context) GetString(key string) string {
	for _, v := range ctx.values {
		if v.Key == key {
			val, ok := v.Val.(string)
			if ok {
				return val
			}
		}
	}
	return ""
}

func (ctx *Context) GetInt(key string) int {
	for _, v := range ctx.values {
		if v.Key == key {
			val, ok := v.Val.(int)
			if ok {
				return val
			}
		}
	}
	return 0
}

func (ctx *Context) GetUint(key string) uint {
	for _, v := range ctx.values {
		if v.Key == key {
			val, ok := v.Val.(uint)
			if ok {
				return val
			}
		}
	}
	return 0
}
