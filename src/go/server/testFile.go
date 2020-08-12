package main

import(
  "syscall/js"
  "strconv"
  "fmt"
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

func setUpBot(this js.Value, i []js.Value) interface{} {
  var strategies string
  space := (" ")[0]
  keyID := js.Global().Get("document").Call("getElementById", i[0].String()).Get("value").String()
  key := js.Global().Get("document").Call("getElementById", i[1].String()).Get("value").String()
  email := js.Global().Get("document").Call("getElementById", i[2].String()).Get("value").String()
  rsi := js.Global().Get("document").Call("getElementById", i[3].String()).Get("value").String()
  macd := js.Global().Get("document").Call("getElementById", i[4].String()).Get("value").String()
  candlesticks := js.Global().Get("document").Call("getElementById", i[5].String()).Get("value").String()
  offset := js.Global().Get("document").Call("getElementById", i[6].String()).Get("value").String()
  info := []string{rsi, macd, candlesticks,offset}
  for _, v := range info {
    if (v[len(v)-1] == space) {
      strategies = strategies + "1"
    } else {
      strategies = strategies + "0"
    }
  }
  fmt.Println(strategies, key, keyID, email)

  return 1
}

func registerCallbacks() {
    js.Global().Set("add", js.FuncOf(add))
    js.Global().Set("subtract", js.FuncOf(subtract))
    js.Global().Set("setUpBot", js.FuncOf(setUpBot))
}

func main() {
    c := make(chan struct{}, 0)

    println("WASM Go Initialized")
    // register functions
    registerCallbacks()
    <-c
}
