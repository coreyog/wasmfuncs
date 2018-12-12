Testing out some wasm helper functions. This project aims to add functions to the global space and check their parameters when called.

Brief example:
```
func main() {
	fmt.Println("Hello, WebAssembly!")

	wasmfuncs.NewFunc("hash", hash, []js.Type{js.TypeString, js.TypeFunction})
	wasmfuncs.NewFunc("print", print, []js.Type{})

	select {}
}

func hash(args []js.Value) {
	output := sha256.Sum256([]byte(args[0].String()))
	args[1].Invoke(fmt.Sprintf("%x", output))
}

func print(args []js.Value) {
	for i := range args {
		fmt.Println(args[i].String())
	}
}

```