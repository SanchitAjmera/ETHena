package main

import (
	"fmt"
  luno "github.com/luno/luno-go"
	"github.com/luno/luno-go/decimal"
)

func testSMA(bot *smaBot) {
	var i int64 = 0
	for i < bot.numOfDecisions {
		bot.tradeSMA()
		i++
	}
}


func testRSI(bot *rsiBot) {
	var i int64 = 0
	for i < bot.numOfDecisions {
		bot.tradeRSI()
		i++
	}
}


var client *luno.Client
var reqPointer *luno.GetTickerRequest

func main() {
	startingFunds := decimal.NewFromInt64(int64(100))

	parseXlsx()

	maxProfit := decimal.NewFromInt64(0)
	maxTradingPeriod := int64(0)

	i:=4

	for i < 121{

		tradingPeriod := int64(i)
		var numOfDecisions int64 = 50000/tradingPeriod
	//	var offset int64 = 40

		pf := portfolio{startingFunds, decimal.NewFromInt64(int64(0)), tradingPeriod/*tradingPeriod*/, tradingPeriod/*tradingPeriod*/, 0}
		//botSMA := smaBot{&pf, decimal.NewFromInt64(offset), numOfDecisions}
		bot := rsiBot{&pf, numOfDecisions, overSold, overBought}
		testRSI(&bot)
		currBid := getBid(bot.pf.currRow)
		portfolioValue := bot.pf.funds.Add(currBid.Mul(bot.pf.stock))
		profit := portfolioValue.Sub(startingFunds)

		if profit.Cmp(maxProfit) == 1 {
				maxProfit = profit
				maxTradingPeriod = bot.pf.tradingPeriod
		}
		i +=2
		}



	days := ((50000 / 60)/ 24)
  fmt.Println("Days: ",days)
//	fmt.Println("Profit/Loss:     £", profit)
  fmt.Println("maximum profit made: £", maxProfit)
	fmt.Println(" at trading periods:  ", maxTradingPeriod)
//	fmt.Println(bot.pf.tradesMade," trades made")

}
