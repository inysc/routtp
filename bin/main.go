package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/inysc/routtp/routtp"
)

func fn(w http.ResponseWriter, r *http.Request) {}

func main() {
	e := routtp.NewNode("", nil)
	e.AddRoute("/evo/rvsl", fn)
	e.AddRoute("/evo/rvsl/", fn)
	e.AddRoute("/evo/rvsl/ecd", fn)
	e.AddRoute("/evo/rvsl/ecdef", fn)
	e.AddRoute("/evo/rvsl/ecd/:a", fn)
	e.AddRoute("/evo/rvsl/ecd/:a/", fn)
	e.AddRoute("/evo/rvsl/ecd/:a/fgsd", fn)
	e.AddRoute("/evo/rvsl/ecd/domain/*", fn)

	routtp.PrintNode(e)

	bs, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile("node.json", bs, 0644)
	if err != nil {
		panic(err)
	}
}
