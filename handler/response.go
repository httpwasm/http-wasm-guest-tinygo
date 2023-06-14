package handler

import (
	"runtime"

	"github.com/httpwasm/http-wasm-guest-tinygo/handler/api"
	"github.com/httpwasm/http-wasm-guest-tinygo/handler/internal/imports"
	"github.com/httpwasm/http-wasm-guest-tinygo/handler/internal/mem"
)

// wasmResponse implements api.Response with imported WebAssembly functions.
type wasmResponse struct{}

// compile-time check to ensure wasmResponse implements api.Response.
var _ api.Response = wasmResponse{}

// GetStatusCode implements the same method as documented on api.Response.
func (r wasmResponse) GetStatusCode() uint32 {
	return imports.GetStatusCode()
}

// SetStatusCode implements the same method as documented on api.Response.
func (r wasmResponse) SetStatusCode(statusCode uint32) {
	imports.SetStatusCode(statusCode)
}

// Headers implements the same method as documented on api.Response.
func (wasmResponse) Headers() api.Header {
	return wasmHeaders
}

// Body implements the same method as documented on api.Response.
func (wasmResponse) Body() api.Body {
	return wasmResponseBody
}

// Trailers implements the same method as documented on api.Response.
func (wasmResponse) Trailers() api.Header {
	return wasmResponseTrailers
}

// GetMethod implements the same template as documented on api.Response.
func (wasmResponse) GetTemplate() string {
	return mem.GetString(imports.GetTemplate)
}

// SetMethod implements the same template as documented on api.Response.
func (wasmResponse) SetTemplate(template string) {
	ptr, size := mem.StringToPtr(template)
	imports.SetTemplate(uintptr(ptr), size)
	runtime.KeepAlive(template) // keep method alive until ptr is no longer needed.
}
