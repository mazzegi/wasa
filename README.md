# wasa
A tiny library to create Web-Apps with Golang and WebAssembly.

## Quick Start
Navigate to `devutil/cmd/server` and start a dev-server by
`go run main.go -bind=:8081  -dist=dist -maingo=../../../example/todomvc/main.go  -wasmexec=../../wasm_exec.js`

In your browser go to http://127.0.0.1:8081/ and look ...
