package main

import (
//	"fmt"
  luno "github.com/luno/luno-go"
	"github.com/luno/luno-go/decimal"
	//"github.com/chenjiandongx/go-echarts/charts"
	//"os"
)

// global variables to retrieve live data
var client *luno.Client
var reqPointer *luno.GetTickerRequest


func main() {

	// processing historical data within excel spreadsheet
	parseXlsx()

  tradingPeriodLR := int64(26)
  tradingPeriodSR := int64(13)

	macdBot := macdBot{
    tradingPeriodLR: 	tradingPeriodLR,
    tradingPeriodSR: 	tradingPeriodSR,
		tradesMade: 		  0,
		numOfDecisions:   0,
    data:             []decimal.Decimal{},
		buyPrice:         decimal.Zero(),
		currRow:				  tradingPeriodLR,
    macdValue:        decimal.Zero(),
	}

  macdBot.initialData()

	for i:= 0; i < int(800);i++ {
		macdBot.trade()

	}
}
