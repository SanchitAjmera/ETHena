package main

import(
  "syscall/js"
  "strconv"
)

// command to recompile lib.wasm: GOARCH=wasm GOOS=js go build -o lib.wasm testFile.go

func add(this js.Value, i []js.Value) interface{} {
    value1 := js.Global().Get("document").Call("getElementById", i[0].String()).Get("value").String()
    value2 := js.Global().Get("document").Call("getElementById", i[1].String()).Get("value").String()

    int1, _ := strconv.Atoi(value1)
    int2, _ := strconv.Atoi(value2)

    js.Global().Get("document").Call("getElementById", i[2].String()).Set("value", int1+int2)
    return int1 + int2
}

func subtract(this js.Value, i []js.Value) interface{} {
    value1 := js.Global().Get("document").Call("getElementById", i[0].String()).Get("value").String()
    value2 := js.Global().Get("document").Call("getElementById", i[1].String()).Get("value").String()

    int1, _ := strconv.Atoi(value1)
    int2, _ := strconv.Atoi(value2)

    js.Global().Get("document").Call("getElementById", i[2].String()).Set("value", int1-int2)
    return int1 - int2
}

func registerCallbacks() {
    js.Global().Set("add", js.FuncOf(add))
    js.Global().Set("subtract", js.FuncOf(subtract))
}

func main() {
    c := make(chan struct{}, 0)

    println("WASM Go Initialized")
    // register functions
    registerCallbacks()
    <-c
}
