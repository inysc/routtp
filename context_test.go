package routtp

import "testing"

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
