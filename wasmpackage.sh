#!/bin/bash
env GOOS=js GOARCH=wasm go build -o space-invaders.wasm space-invaders

cp $(go env GOROOT)/misc/wasm/wasm_exec.js .

