package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/inysc/routtp/routtp"
)

func fn(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("<goid:%d>%s %s\n", runtime.Goid(), r.Method, r.RequestURI)
	// fmt.Printf("%s\n", debug.Stack())
	ctx := r.Context().(*routtp.Context)
	fmt.Printf("<goid:%d>ctx.Param:<%+v>\n", runtime.Goid(), ctx.Param)
}

func main() {
	log.SetFlags(log.Lshortfile)
	gin.SetMode(gin.ReleaseMode)

	router := routtp.New()
	group := router.Group(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("group\n")
	})
	router.Use(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context().(*routtp.Context)
		fmt.Printf("param<%+v>\n", ctx.Param)
		fmt.Println("in middleware1<next前>")
		ctx.Next()
		fmt.Println("in middleware1<next>后")
	}, func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("in middleware 2")
	}, func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context().(*routtp.Context)
		ctx.Abort()
		fmt.Println("in middleware 3<abort>")
		// fmt.Println("in middleware 3<next前>")
		// ctx.Next()
		// fmt.Println("in middleware 3<next>后")
	}, func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context().(*routtp.Context)
		fmt.Println("in middleware 4<next前>")
		ctx.Next()
		fmt.Println("in middleware 4<next>后")
	})
	router.GET("/evo/rvsl", fn)

	router.GET("/evo/rvsl/", fn)
	group.GET("/evo/rvsl/ecd", fn)
	router.GET("/evo/rvsl/ecdef", fn)
	router.GET("/evo/rvsl/ecdeg", fn)
	router.GET("/evo/rvsl/ecd/:a", fn)
	router.GET("/evo/rvsl/ecd/:a/", fn)
	router.GET("/evo/rvsl/ecd/:a/fgd", fn)
	router.GET("/evo/rvsl/ecd/domain/*", fn)

	routtp.PrintNode(router.Method[0].Val)

	println("start :8080")
	http.ListenAndServe(":8080", router)
}
