// +build js,wasm

package wasmfuncs

import (
	"fmt"
	"strings"
	"syscall/js"
)

// NewFunc adds a function to the global object
func NewFunc(funcName string, fn func([]js.Value), signature []js.Type) {
	js.Global().Set(funcName, js.NewCallback(func(args []js.Value) {
		if len(signature) != 0 && !checkArgs(args, signature) {
			fmt.Println("bad WASM function call")
			fmt.Printf("  expected: %s(%s)\n", funcName, typesToStrings(signature))
			fmt.Printf("    actual: %s(%s)\n", funcName, typesToStrings(valuesToTypes(args)))
			return
		}
		fn(args)
	}))

	wasm := js.Global().Get("wasm")
	if wasm.Type() == js.TypeUndefined {
		wasm = js.ValueOf(make(map[string]interface{}))
	}
	registeredFunctions := wasm.Get("registeredFunctions")
	if registeredFunctions.Type() == js.TypeUndefined {
		registeredFunctions = js.ValueOf(make(map[string]interface{}))
	}
	registeredFunctions.Set(funcName, fmt.Sprintf("(%s)", typesToStrings(signature)))
	wasm.Set("registeredFunctions", registeredFunctions)
	js.Global().Set("wasm", wasm)
	// let wasm = {
	// 	registeredFunctions: {
	// 		func1: sig1,
	// 		func2: sig2,
	// 		...
	// 	}
	// }
}

func checkArgs(args []js.Value, types []js.Type) bool {
	if len(args) != len(types) {
		return false
	}
	for i := range args {
		if args[i].Type() != types[i] {
			return false
		}
	}
	return true
}

func valuesToTypes(values []js.Value) (types []js.Type) {
	types = make([]js.Type, len(values))
	for i := range values {
		types[i] = values[i].Type()
	}
	return types
}

func typesToStrings(types []js.Type) (output string) {
	arr := make([]string, len(types))
	for i := range types {
		switch types[i] {
		case js.TypeBoolean:
			arr[i] = "bool"
		case js.TypeFunction:
			arr[i] = "function"
		case js.TypeNull:
			arr[i] = "null"
		case js.TypeNumber:
			arr[i] = "number"
		case js.TypeObject: // JS arrays are objects?!?! WTF
			arr[i] = "object"
		case js.TypeString:
			arr[i] = "string"
		case js.TypeSymbol:
			arr[i] = "symbol"
		case js.TypeUndefined:
			arr[i] = "undefined"
		}
	}
	return strings.Join(arr, ", ")
}
