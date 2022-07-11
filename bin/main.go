package main

import (
	"fmt"
	"net/http"
	"runtime"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/inysc/routtp/routtp"
)

func fn(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s %s<%d>\n", r.Method, r.RequestURI, runtime.Goid())
	fmt.Printf("%s\n", debug.Stack())
	ctx := r.Context().(*routtp.Context)
	for _, v := range ctx.Param {
		fmt.Printf("key<%s>, value<%s>\n", v.Key, v.Val)
	}
	println("\n")
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	router := routtp.New()
	router.GET("/evo/rvsl", fn)
	router.GET("/evo/rvsl/", fn)
	router.GET("/evo/rvsl/ecd", fn)
	router.GET("/evo/rvsl/ecdef", fn)
	router.GET("/evo/rvsl/ecd/:a", fn)
	router.GET("/evo/rvsl/ecd/:a/", fn)
	router.GET("/evo/rvsl/ecd/:a/fgd", fn)
	router.GET("/evo/rvsl/ecd/domain/*", fn)

	routtp.PrintNode(router.Method[0].Val)

	println("start :8080")
	// http.ListenAndServe(":8080", router)
}
