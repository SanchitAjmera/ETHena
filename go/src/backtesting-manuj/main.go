package main

import "github.com/luno/luno-go/decimal"

func test(bot smaBot) {
	var i int64 = 0
	for i < bot.numOfDecisions {
		bot.trade()
		i++
	}
}

func main() {
	var numOfDecisions int64 = 20
	var offset int64 = 127
	pf := portfolio{decimal.NewFromInt64(int64(100)), decimal.NewFromInt64(int64(0)), int64(30), int64(35)}
	bot := smaBot{pf, decimal.NewFromInt64(offset), numOfDecisions}
	test(bot)
}
