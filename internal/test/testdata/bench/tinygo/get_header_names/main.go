package main

import (
	"github.com/httpwasm/http-wasm-guest-tinygo/handler"
	"github.com/httpwasm/http-wasm-guest-tinygo/handler/api"
)

func main() {
	handler.HandleRequestFn = GetHeaderNames
}

func GetHeaderNames(req api.Request, resp api.Response) (next bool, reqCtx uint32) {
	_ = req.Headers().Names()
	return // this is a benchmark, so skip the next handler.
}
