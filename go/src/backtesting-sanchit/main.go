package main

import (
	"fmt"
	"github.com/luno/luno-go/decimal"
)

func test(bot *smaBot) {
	var i int64 = 0
	for i < bot.numOfDecisions {
		bot.trade()
		i++
	}
}

func main() {

	startingFunds := decimal.NewFromInt64(int64(100))

	parseXlsx()

	periods := [12]int64{5, 15,20,25, 30,35,40,45,50,55, 60, 240}
	var bestPeriod int64 = 0
	var bestOffset int64 = 0
	var maxProfit decimal.Decimal = decimal.Zero()
	for _, v := range periods {
		var tradingPeriod int64 = int64(v)
		var numOfDecisions int64 = 50000/tradingPeriod
		days := (numOfDecisions * tradingPeriod) / (60 * 24)

		for j := 5; j < 200; j += 5 {
			var offset int64 = int64(j)
			pf := portfolio{startingFunds, decimal.NewFromInt64(int64(0)), tradingPeriod, int64(tradingPeriod), 0}
			bot := smaBot{pf, decimal.NewFromInt64(offset), numOfDecisions}
			test(&bot)

			currBid := getBid(bot.pf.currRow)
			portfolioValue := bot.pf.funds.Add(currBid.Mul(bot.pf.stock))
			profit := portfolioValue.Sub(startingFunds)

			if (profit.Cmp(maxProfit) == 1) {
				maxProfit = profit
				bestPeriod = tradingPeriod
				bestOffset = offset
			}

			fmt.Println("Days: ",days,"  Trading Period: ",tradingPeriod,"  Offset: ",offset)
			fmt.Println("Profit/Loss:     £", profit)
			fmt.Println(bot.pf.tradesMade," trades made")
		}
	}
	fmt.Println("Max Profit:    £", maxProfit)
	fmt.Println("Best Period:    ", bestPeriod)
	fmt.Println("Best Offset:    ", bestOffset)
}
