package main

import(
  "fmt"
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

func Run(this js.Value, i []js.Value) interface{} {
  keyID := js.Global().Get("document").Call("getElementById", i[0].String()).Get("value").String()
  key := js.Global().Get("document").Call("getElementById", i[1].String()).Get("value").String()
  email := js.Global().Get("document").Call("getElementById", i[2].String()).Get("value").String()
  rsi := js.Global().Get("document").Call("getElementById", i[3].String()).Get("value").String()
  macd := js.Global().Get("document").Call("getElementById", i[4].String()).Get("value").String()
  candlesticks := js.Global().Get("document").Call("getElementById", i[5].String()).Get("value").String()
  js.Global().Get("document").Call("getElementById", i[6].String()).Set("value", 2)
  fmt.Println(keyID, key, email, rsi, macd, candlesticks)
  return 1
}

func registerCallbacks() {
    js.Global().Set("add", js.FuncOf(add))
    js.Global().Set("subtract", js.FuncOf(subtract))
    js.Global().Set("Run", js.FuncOf(Run))
}

func main() {
    c := make(chan struct{}, 0)

    println("WASM Go Initialized")
    // register functions
    registerCallbacks()
    <-c
}
