package main

import (
	"fmt"
  luno "github.com/luno/luno-go"
	"github.com/luno/luno-go/decimal"
)

func populatePastBids () {
	//Populating past bids with 1 tradingPeriod worth of data
	for i < b.tradingPeriod {
		bot.pastBids[i] = getCurrBid()
		fmt.Println("Filling past bids: ",i,"/",b.tradingPeriod)
		time.Sleep(time.Minute)
	}
}

// test function for the RSI bot
func test(bot *rsiBot) {
	var i int64 = 0
	for {
		bot.trade()
	}
}

// global variables to retrieve live data
var client *luno.Client
var reqPointer *luno.GetTickerRequest
var pair string

func main() {

	pair := "XBTGBP"
	client, reqPointer = getTickerRequest()

	accountReq := CreateAccountRequest{Currency: "GBP", Name: "rsiAccount"}
	account, err := client.CreateAccount(context.Background(), &accountReq)

	if err != nil {panic(err)}


	// initial funds at the start of the trade period
	startingFunds := getAsset("GBP")

	// initialising values within bot portfolio
	tradingPeriod := int64(20)
	numOfDecisions := int64(50000/tradingPeriod)
	stopLossMultDecimal := decimal.NewFromString("0.97")

	rsiLowerLim := decimal.NewFromInt64(30)
	rsiUpperLim := decimal.NewFromInt64(70)

	// initialising bot
	bot := rsiBot{
		 	tradingPeriod: 	tradingPeriod,
			tradesMade: 		0,
			numOfDecisions: 0,
			stopLoss: 			decimal.Zero(),
			stopLossMult: 	stopLossMultDecimal,
			pastBids: 			[tradingPeriod]decimal.Decimal{},
			overSold: 			rsiLowerLim,
			overBought: 		rsiUpperLim
	}


	test(&bot)

	// calculating overall profit
	currBid := getCurrBid()
	gbpBalance := getAsset("GBP")
	xbtBalance := getAsset("XBT")
	portfolioValue := gbpBalance.Add(currBid.Mul(xbtBalance))
	profit := portfolioValue.Sub(startingFunds)

}
