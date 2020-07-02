package main

import (
	"fmt"
  luno "github.com/luno/luno-go"
	"github.com/luno/luno-go/decimal"
)

// test function for the SMA bot
func testSMA(bot *smaBot) {
	var i int64 = 0
	for i < bot.numOfDecisions {
		bot.tradeSMA()
		i++
	}
}


// test function for the RSI bot
func testRSI(bot *rsiBot) {
	var i int64 = 0
	for i < bot.numOfDecisions {
		bot.tradeRSI()
		i++
	}
}

// global variables to retrieve live data
var client *luno.Client
var reqPointer *luno.GetTickerRequest


func main() {
	// initial funds at the start of the trade period
	startingFunds := decimal.NewFromInt64(int64(100))
	// processing historical data within excel spreadsheet
	parseXlsx()
	// initialising max values
	maxProfit := decimal.NewFromInt64(0)
	maxTradingPeriod := int64(0)
	maxOverSold := int64(0)
	maxOverBought := int64(0)

	// iterating through different metrics for the bot
	for x := 0; x < 50; x+=10 {
		for y := 100; y > 50; y-=10 {
			i:=4
			for i < 60{
				// initialising values within bot portfolio
				overSold := int64(x)
				overBought := int64(y)
				tradingPeriod := int64(i)
				var numOfDecisions int64 = 50000/tradingPeriod
				//	var offset int64 = 40
				// initialising portfolio
				pf := portfolio{startingFunds, decimal.NewFromInt64(int64(0)), tradingPeriod/*tradingPeriod*/, tradingPeriod/*tradingPeriod*/, 0}
				//botSMA := smaBot{&pf, decimal.NewFromInt64(offset), numOfDecisions}
				// initialising bot
				bot := rsiBot{&pf, numOfDecisions, overSold, overBought}
				testRSI(&bot)
				// calculating overall profit
				currBid := getBid(bot.pf.currRow)
				portfolioValue := bot.pf.funds.Add(currBid.Mul(bot.pf.stock))
				profit := portfolioValue.Sub(startingFunds)
				// check to see if a greater max profit was generated
				if profit.Cmp(maxProfit) == 1 {
					maxProfit = profit
					maxOverSold = bot.overSold
					maxOverBought = bot.overBought
					maxTradingPeriod = bot.pf.tradingPeriod
					fmt.Println("maximum profit made: £", maxProfit)
					fmt.Println("at trading periods:  ", maxTradingPeriod)
					fmt.Println("upper bound: 			  ", maxOverBought)
					fmt.Println("upper Sold:  				", maxOverSold)
					fmt.Println(".")
				}
				i +=4
			}
		}
	}
	// printing out results
	days := ((50000 / 60)/ 24)
  fmt.Println("Days: ",days)
//	fmt.Println("Profit/Loss:     £", profit)
  fmt.Println("maximum profit made: £", maxProfit)
	fmt.Println(" at trading periods:  ", maxTradingPeriod)
	fmt.Println(" 			 upper bound:  ", maxOverBought)
	fmt.Println(" 			 upper Sold:  ", maxOverSold)
//	fmt.Println(bot.pf.tradesMade," trades made")

}
