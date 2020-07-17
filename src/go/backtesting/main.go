package main

import (
//	"fmt"
  luno "github.com/luno/luno-go"
	"github.com/luno/luno-go/decimal"
	//"github.com/chenjiandongx/go-echarts/charts"
	//"os"
)

// test function for the SMA bot
func testSMA(bot *smaBot) {
	var i int64 = 0
	for i < bot.NumOfDecisions {
		bot.tradeSMA()
		i++
	}
}

// test function for the RSI bot
func testRSI(bot *rsiBot) {
	var i int64 = 0
	for i < bot.NumOfDecisions {
		bot.tradeRSI()
		i++
	}
}


func testMACD(bot *macdBot) {
	var i int64 = 0
	for i < bot.NumOfDecisions {
		bot.tradeMACD()
		i++
	}
}

// global variables to retrieve live data
var client *luno.Client
var reqPointer *luno.GetTickerRequest


func main() {

	// processing historical data within excel spreadsheet
	parseXlsx()

	TradingPeriod := int64(20)

	candle := candleBot{
		TradingPeriod: 	TradingPeriod, //180 minute candlestick
		TradesMade: 		0,
		NumOfDecisions: 0,
		queue: 					[]candlestick{},
		ReadyToBuy: 		true,
		BuyPrice:       decimal.Zero(),
		currRow:				1,
	}

	candle.fillQueue(3)

	for i:= 0; i < int(800/TradingPeriod);i++ {
		candle.trade()

	}
}
