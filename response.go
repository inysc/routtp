package routtp

import (
	"encoding/json"
	"net/http"
)

func (ctx *Context) write(statusCode int, contentType string, body []byte) {
	ctx.WriteHeader(statusCode)
	ctx.Write(body)
}

type JSON = map[string]any

func (ctx *Context) JSON(statusCode int, body any) {
	bs, err := json.Marshal(body)
	if err != nil {
		ctx.Exception(&exception{
			code:   http.StatusInternalServerError,
			status: http.StatusInternalServerError,
			msg:    "faild to marshal response body",
			data:   "",
		})
		return
	}
	ctx.write(statusCode, "application/json", bs)
}

// 成功响应
func (ctx *Context) Success(body any) {
	switch bd := body.(type) {
	case string:
		ctx.write(200, "text/plain", StringToBytes(bd))
	case []byte:
		ctx.write(200, "text/plain", bd)
	default:
		bs, err := json.Marshal(body)
		if err != nil {
			ctx.Exception(&exception{
				code:   http.StatusInternalServerError,
				status: http.StatusInternalServerError,
				msg:    "faild to marshal response body",
				data:   "",
			})
			return
		}
		ctx.write(200, "application/json", bs)
	}
}

// 异常响应
func (ctx *Context) Exception(resp Response) {
	ctx.JSON(resp.Status(), JSON{
		"code": resp.Code(),
		"msg":  resp.Msg(),
		"data": resp.Data(),
	})
}

type Response interface {
	Code() int
	Msg() string
	Data() any
	Status() int
}

type exception struct {
	code   int
	status int
	msg    string
	data   string
}

func (ep *exception) Code() int {
	return ep.code
}

func (ep *exception) Msg() string {
	return ep.msg
}

func (ep *exception) Data() any {
	return ep.data
}

func (ep *exception) Status() int {
	return ep.status
}
