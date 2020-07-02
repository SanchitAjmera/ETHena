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

	// initialising values within bot portfolio
	tradingPeriod := int64(20)
	var numOfDecisions int64 = 50000/tradingPeriod
	stopLossMultiplier := decimal.NewFromInt64(97)
	stopLossMultDecimal := stopLossMultiplier.DivInt64(int64(100))
	// var offset int64 = 40
	// initialising portfolio
	pf := portfolio{funds: startingFunds,
									stock: decimal.Zero(),
									tradingPeriod: tradingPeriod,
									currRow: tradingPeriod,
									tradesMade: 0,
									stopLoss: decimal.Zero(),
									stopLossMult: stopLossMultDecimal}

	// bot := smaBot{&pf, decimal.NewFromInt64(offset), numOfDecisions}
	// initialising bot
	bot := rsiBot{&pf, numOfDecisions, int64(30), int64(70)}
	testRSI(&bot)
	// calculating overall profit
	currBid := getBid(bot.pf.currRow)
	portfolioValue := bot.pf.funds.Add(currBid.Mul(bot.pf.stock))
	profit := portfolioValue.Sub(startingFunds)
	/*
	// initialising max values
	maxProfit := decimal.NewFromInt64(0)
	maxStopLossMult := decimal.Zero()
	// iterating through different metrics for the bot
	for x := 75; x < 100; x+=1 {
		// initialising values within bot portfolio
		tradingPeriod := int64(20)
		var numOfDecisions int64 = 50000/tradingPeriod
		stopLossMultiplier := decimal.NewFromInt64(int64(x))
		stopLossMultDecimal := stopLossMultiplier.Div(decimal.NewFromInt64(100),8)
		//	var offset int64 = 40
		// initialising portfolio
		pf := portfolio{startingFunds, decimal.NewFromInt64(int64(0)), tradingPeriod, tradingPeriod, 0, decimal.Zero(), stopLossMultDecimal}
		//botSMA := smaBot{&pf, decimal.NewFromInt64(offset), numOfDecisions}
		// initialising bot
		bot := rsiBot{&pf, numOfDecisions, int64(30), int64(70)}
		testRSI(&bot)
		// calculating overall profit
		currBid := getBid(bot.pf.currRow)
		portfolioValue := bot.pf.funds.Add(currBid.Mul(bot.pf.stock))
		profit := portfolioValue.Sub(startingFunds)
		// check to see if a greater max profit was generated
		if profit.Cmp(maxProfit) == 1 {
			maxProfit = profit
			maxStopLossMult = bot.pf.stopLossMult
			//fmt.Println("maximum profit made: £", maxProfit)
			//fmt.Println("at trading periods:  ", maxTradingPeriod)
			//fmt.Println("upper bound: 			  ", maxOverBought)
			//fmt.Println("upper Sold:  				", maxOverSold)
			//fmt.Println(".")
		}
	}
	*/
	// printing out results
	days := ((50000 / 60)/ 24)
  fmt.Println("Days: ",days)
	fmt.Println("Profit/Loss:     £", profit)
//  fmt.Println("maximum profit made: £", maxProfit)
//	fmt.Println(" at trading periods:  ", maxTradingPeriod)
//	fmt.Println("stop loss multiplier: ", maxStopLossMult)
//	fmt.Println(" 			 upper Sold:  ", maxOverSold)
	fmt.Println(bot.pf.tradesMade," trades made")

}
