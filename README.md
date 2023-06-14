[![Build](https://github.com/httpwasm/http-wasm-guest-tinygo/workflows/build/badge.svg)](https://github.com/httpwasm/http-wasm-guest-tinygo)
[![Go Report Card](https://goreportcard.com/badge/github.com/httpwasm/http-wasm-guest-tinygo)](https://goreportcard.com/report/github.com/httpwasm/http-wasm-guest-tinygo)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)

# http-wasm Guest Library for TinyGo

[http-wasm][1] is HTTP client middleware implemented in [WebAssembly][2].
This is a [TinyGo WASI][3] library that implements the [Guest ABI][4].

## Example
The following is an [example](examples/router) of routing middleware:

```go
package main

import (
	"strings"

	"github.com/httpwasm/http-wasm-guest-tinygo/handler"
	"github.com/httpwasm/http-wasm-guest-tinygo/handler/api"
)

func main() {
	handler.HandleRequestFn = handleRequest
}

// handle implements a simple HTTP router.
func handleRequest(req api.Request, resp api.Response) (next bool, reqCtx uint32) {
	// If the URI starts with /host, trim it and dispatch to the next handler.
	if uri := req.GetURI(); strings.HasPrefix(uri, "/host") {
		req.SetURI(uri[5:])
		next = true // proceed to the next handler on the host.
		return
	}

	// Serve a static response
	resp.Headers().Set("Content-Type", "text/plain")
	resp.Body().WriteString("hello")
	return // skip any handlers as the response is written.
}
```

If you make changes, you can rebuild this with TinyGo v0.28 or higher like so:
```sh
tinygo build -o examples/router/main.wasm -scheduler=none --no-debug -target=wasi examples/router/main.go
```

There are also more [examples](examples) you may wish to try out!

# WARNING: This is an early draft

The current maturity phase is early draft. Once this is integrated with
[coraza][5] and [dapr][6], we can begin discussions about compatability.

[1]: https://github.com/http-wasm
[2]: https://webassembly.org/
[3]: https://wazero.io/languages/tinygo/
[4]: https://github.com/httpwasm/http-wasm-abi
[5]: https://github.com/corazawaf/coraza-proxy-wasm
[6]: https://github.com/httpwasm/components-contrib/

#  个人理解
## wasm/guest处理流程  

- 使用```//go:export```注解导出函数,httpwasm目前导出了两个函数 ```handle_request``` 和 ```handle_response```,  
```handleRequest```  调用 ```HandleRequestFn```  
```handle_response``` 调用  ```HandleResponseFn```  
```go
// 代码在 ./handler/handler.go#20

// handleRequest is only exported to the host.
// wasm导出的函数,宿主host可以调用,使用 go:export 标记
//
//go:export handle_request
func handleRequest() (ctxNext uint64) { //nolint
	next, reqCtx := HandleRequestFn(wasmRequest{}, wasmResponse{})
}


// handleResponse is only exported to the host.
// wasm导出的函数,宿主host可以调用,使用 go:export 标记
//
//go:export handle_response
func handleResponse(reqCtx uint32, isError uint32) { //nolint
	HandleResponseFn(reqCtx, wasmRequest{}, wasmResponse{}, isErrorB)
}
```
- 在wasm的main函数中指定```HandleRequestFn``` 和```HandleResponseFn```的实现. host-->guest.handle_request-->HandleRequestFn-->具体实现
```go
func main() {
	handler.HandleRequestFn = handleRequest
}

// handleRequest implements a simple HTTP router.
func handleRequest(req api.Request, resp api.Response) (next bool, reqCtx uint32) {
}
```
```HandleRequestFn``` 有 ```api.Request``` 和 ```api.Response``` 两个接口参数,用来绑定host的request和reponse,下面说下如何实现
```api.Request```接口实现是```wasmRequest```,```api.Response```接口实现是```wasmResponse```.实现再调用host 导入的函数,根据指定的数据格式进行交换,所以 wasm 默认不能处理并发请求.  

- 导入host函数
复杂的逻辑可以放到host端完成,通过imports导入函数,由guest wasm进行调用host的实现,主要通过指针地址和数据长度,完成数据交换  
所有的函数名均为 小写 下划线进行连接,详见:https://http-wasm.io/rationale/#why-is-everything-lower_snake_case-instead-of-lower-hyphen-case  
```./handler/internal/imports/imports.go``` 中,使用```//go:build tinygo.wasm```标识tinygo编译,以本次新增的方法为例

```go
//go:wasmimport http_handler get_template
func getTemplate(ptr uint32, limit BufLimit) (len uint32)

//go:wasmimport http_handler set_template
func setTemplate(ptr, size uint32)
```
  
``` //go:wasmimport http_handler get_template ``` wasm的模块名称是```http_handler```,host加载wasm的二进制文件,初始化模块名称时需要指定为```http_handler```,host的import导入函数,名称为 ```get_template```


