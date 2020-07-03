package main

import (
	"fmt"
  luno "github.com/luno/luno-go"
	"github.com/luno/luno-go/decimal"
	//"github.com/chenjiandongx/go-echarts/charts"
	//"os"
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


func testMACD(bot *macdBot) {
	var i int64 = 0
	for i < bot.numOfDecisions {
		bot.tradeMACD()
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
/*
	// initialising values within bot portfolio
	tradingPeriod := int64(1)
	currRow := int64(40)
	results := make([]decimal.Decimal, 0)
	var numOfDecisions int64 = 50000/tradingPeriod
	stopLossMultiplier := decimal.NewFromInt64(97)
	stopLossMultDecimal := stopLossMultiplier.Div(decimal.NewFromInt64(100),8)
	// var offset int64 = 40
	// initialising portfolio
	pf := portfolio{startingFunds, decimal.NewFromInt64(int64(0)), tradingPeriod, currRow, 0, decimal.Zero(), stopLossMultDecimal}
	// bot := smaBot{&pf, decimal.NewFromInt64(offset), numOfDecisions}
	// initialising bot
	bot := macdBot{&pf, numOfDecisions, int64(10), currRow, results, decimal.Zero()}
	testMACD(&bot)
	// calculating overall profit
	currBid := getBid(bot.pf.currRow)
	portfolioValue := bot.pf.funds.Add(currBid.Mul(bot.pf.stock))
	profit := portfolioValue.Sub(startingFunds)
*/
	// initialising max values
	maxProfit := decimal.NewFromInt64(0)
	maxSRDuration := 0
	maxLRDuration := 0
	// iterating through different metrics for the bot
	for sr := 5 ; sr < 10; sr+=5 {
		for lr := sr+5;  lr< sr+10; lr+=10 {

			// initialising values within bot portfolio
			// initialising values within bot portfolio
			tradingPeriod := int64(1)

			results := make([]decimal.Decimal, 0)
			var numOfDecisions int64 = 50000/tradingPeriod
			stopLossMultiplier := decimal.NewFromInt64(97)
			stopLossMultDecimal := stopLossMultiplier.Div(decimal.NewFromInt64(100),8)
			// var offset int64 = 40
			// initialising portfolio
			pf := portfolio{startingFunds, decimal.NewFromInt64(int64(0)), tradingPeriod, int64(300), 0, decimal.Zero(), stopLossMultDecimal}
			// bot := smaBot{&pf, decimal.NewFromInt64(offset), numOfDecisions}
			// initialising bot
			bot := macdBot{&pf, numOfDecisions, int64(30), int64(300), results, decimal.Zero()}
			testMACD(&bot)
			// calculating overall profit
			currBid := getBid(bot.pf.currRow)
			portfolioValue := bot.pf.funds.Add(currBid.Mul(bot.pf.stock))
			profit := portfolioValue.Sub(startingFunds)

			if profit.Cmp(maxProfit) == 1 {
				maxProfit = profit
				maxSRDuration = 30
				maxLRDuration = 300
				fmt.Println("max Profit" , maxProfit)
				fmt.Println("max SR" , maxSRDuration)
				fmt.Println("max LR" , maxLRDuration, "\n")
			}
		}
	}

	// printing out results
	//days := ((50000 / 60)/ 24)
  //fmt.Println("Days: ",days)
	//fmt.Println("Profit/Loss:     £", profit)
	fmt.Println("max Profit" , maxProfit)
	fmt.Println("max SR" , maxSRDuration)
	fmt.Println("max LR" , maxLRDuration, "\n")
//	fmt.Println(bot.pf.tradesMade," trades made")
//	xAxis := make([]decimal.Decimal,0)
//	for i := 0;  i < len(bot.diffs); i++ {
//		xAxis = append(xAxis, decimal.NewFromInt64(int64(i)))
//	}

//	bar := charts.NewBar()
//	bar.AddXAxis(xAxis).AddYAxis("diffs", bot.diffs)

	//f, err := os.Create("bar.html")
	//if err != nil {panic(err)}
	//bar.Render(f)


}
