package routtp

import (
	"os"
	"testing"
)

func TestPrefix(t *testing.T) {
	ctx := &Context{}
	for _, v := range []struct {
		A string
		B string
	}{
		{A: "/rvsl/:ecd", B: "/rvsl//"},         // ecd<> // 无法正确匹配
		{A: "/rvsl/:ecd", B: "/rvsl/asdf"},      // ecd<asdf>
		{A: "/rvsl/:ecd", B: "/rvsl/asdf/ghjk"}, // ecd<asdf>
		{A: "/rvsl/*ecd", B: "/rvsl/asdf/ghjk"},
	} {
		i, j := ctx.prefix(v.A, v.B)
		t.Logf("i<%d>, j<%d>, param<%+v>", i, j, ctx.param)
		ctx.clean()
	}
}

func TestNodeFprint(t *testing.T) {
	n := &Node{Path: "/"}

	{
		n.AddRoute("rvsl/:ecd", func(ctx *Context) {})
		n.AddRoute("rvsl/:ecd", func(ctx *Context) {}, func(ctx *Context) {})
		n.AddRoute("rvsl/:ecd/:fg", func(ctx *Context) {}, func(ctx *Context) {})
		n.AddRoute("rvsl/:ecd/afgg", func(ctx *Context) {}, func(ctx *Context) {})
		n.AddRoute("rvsl/:ecd/afgh", func(ctx *Context) {}, func(ctx *Context) {})
		n.AddRoute("rvsl/:ecd/afgi", func(ctx *Context) {}, func(ctx *Context) {})
		n.AddRoute("rvsl/:ecd/:fgh/tija", func(ctx *Context) {}, func(ctx *Context) {})
		n.AddRoute("rvsl/:ecd/:fgh/tijb", func(ctx *Context) {}, func(ctx *Context) {})
		n.AddRoute("rvsl/:ecd/:fgh/tijc", func(ctx *Context) {}, func(ctx *Context) {})
		n.AddRoute("rvsl/:ecd/:fgh/tijd", func(ctx *Context) {}, func(ctx *Context) {})
		n.AddRoute("rvsl/:ecd/:fgh/tije", func(ctx *Context) {}, func(ctx *Context) {})
		n.AddRoute("rvsl/:ecd/:fgh/tijf", func(ctx *Context) {}, func(ctx *Context) {})
		n.AddRoute("rvsl/:ecd/:fgh/tijg", func(ctx *Context) {}, func(ctx *Context) {})
		n.AddRoute("rvsl/:ecd/:fgh/tijh", func(ctx *Context) {}, func(ctx *Context) {})
	}

	n.Fprint(os.Stdout, 0)
}
