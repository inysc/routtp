package routtp

import (
	"encoding/json"
	"net/http"
)

func (ctx *Context) write(statusCode int, contentType string, body []byte) {
	ctx.Response.WriteHeader(statusCode)
	ctx.Response.Write(body)
}

func (ctx *Context) STRING(statusCode int, body string) {
	ctx.write(statusCode, "text/plain", StringToBytes(body))
}

type JSON = map[string]any

func (ctx *Context) JSON(statusCode int, body JSON) {
	bs, err := json.Marshal(body)
	if err != nil {
		ctx.write(http.StatusInternalServerError, "text/plain", StringToBytes(err.Error()))
		return
	}
	ctx.write(statusCode, "application/json", bs)
}
