package main

import (
	"fmt"
  luno "github.com/luno/luno-go"
	"github.com/luno/luno-go/decimal"
)

func test(bot *smaBot) {
	var i int64 = 0
	for i < bot.numOfDecisions {
		bot.trade()
		i++
	}
}

var client *luno.Client
var reqPointer *luno.GetTickerRequest

func main() {
	startingFunds := decimal.NewFromInt64(int64(100))

	parseXlsx()

	var tradingPeriod int64 = 20
	var numOfDecisions int64 = 50000/tradingPeriod
	var offset int64 = 40

	pf := portfolio{startingFunds, decimal.NewFromInt64(int64(0)), tradingPeriod, int64(tradingPeriod), 0}
	bot := smaBot{pf, decimal.NewFromInt64(offset), numOfDecisions}
	test(&bot)

	currBid := getBid(bot.pf.currRow)
	portfolioValue := bot.pf.funds.Add(currBid.Mul(bot.pf.stock))
	profit := portfolioValue.Sub(startingFunds)

	days := (numOfDecisions * tradingPeriod) / (60 * 24)
	fmt.Println("Days: ",days,"  Trading Period: ",tradingPeriod,"  Offset: ",offset)
	fmt.Println("Profit/Loss:     £", profit)
	fmt.Println(bot.pf.tradesMade," trades made")
}
