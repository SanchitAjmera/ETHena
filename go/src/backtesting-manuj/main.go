package main

import (
        "github.com/luno/luno-go/decimal"
)

func test (bot SMA_bot) {
  for i := 0; i < bot.numOfDecisions, i++ {
    bot.trade(bot.pf)
  }
}

func main () {
  numOfDecisions := 20
  offset := 127
  pf := portfolio{100, 0, 30, 35}
  bot := SMA_bot{pf, offset, numOfDecisions}
  test(bot)
}
