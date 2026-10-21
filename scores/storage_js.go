//go:build js && wasm

package score

import (
	"fmt"
	"strconv"
	"syscall/js"
)

const storageKey = "space-invaders.highscore.v1"

var session int
var unavailable bool

func storage(action func(js.Value)) (ok bool) {
	defer func() {
		if recover() != nil {
			ok = false
		}
	}()
	// Reflect.get translates a throwing property getter into a recoverable panic.
	value := js.Global().Get("Reflect").Call("get", js.Global(), "localStorage")
	if value.IsNull() || value.IsUndefined() {
		return false
	}
	action(value)
	return true
}

func fallback() {
	unavailable = true
	fmt.Println("Browser high-score storage unavailable; using this session only.")
}

func Load() int {
	if unavailable {
		return session
	}
	if !storage(func(value js.Value) {
		data := value.Call("getItem", storageKey)
		if !data.IsNull() && !data.IsUndefined() {
			session = parse(data.String())
		}
	}) {
		fallback()
	}
	return session
}

func Save(value int) error {
	session = max(0, value)
	if !unavailable && !storage(func(value js.Value) { value.Call("setItem", storageKey, strconv.Itoa(session)) }) {
		fallback()
	}
	return nil
}
