//go:build js && wasm

package platform

import "syscall/js"

const Browser = true

func Ready() {
	callback := js.Global().Get("spaceInvadersReady")
	if callback.Type() == js.TypeFunction {
		callback.Invoke()
	}
}
