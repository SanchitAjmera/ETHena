package main

import (
	// "fmt"
	// "context"
	// "time"
  luno "github.com/luno/luno-go"
	// "github.com/luno/luno-go/decimal"
)

/*TODO: Change sleep from second to minute*/



// test function for the RSI bot
func test(b *rsiBot) {
	populatePastAsks(b)
	for {

		// calculating overall profit
		b.trade()

		/*
		gbpBalance := getAsset("GBP")
		xbtBalance := getAsset("XBT")
		portfolioValue := gbpBalance.Add(getCurrBid().Mul(xbtBalance))

		fmt.Println("Portfolio Value: £", portfolioValue)
		fmt.Println("")
		*/

	}
}

// Global Variables
var client *luno.Client
var reqPointer *luno.GetTickerRequest
var pair string

func main() {

	pair = "XBTGBP"
	client, reqPointer = getTickerRequest()


	//accountReq := luno.CreateAccountRequest{Currency: "GBP", Name: "rsiAccount"}
	//_, err := client.CreateAccount(context.Background(), &accountReq)

	//if err != nil {panic(err)}

	// initialising values within bot portfolio
	/*
	tradingPeriod := int64(20)
	stopLossMultDecimal := decimal.NewFromFloat64(0.97, 8)

	rsiLowerLim := decimal.NewFromInt64(30)
	rsiUpperLim := decimal.NewFromInt64(70)

	// initialising bot
	bot := rsiBot{
		 	tradingPeriod: 	tradingPeriod,
			tradesMade: 		0,
			numOfDecisions: 0,
			stopLoss: 			decimal.Zero(),
			stopLossMult: 	stopLossMultDecimal,
			pastAsks: 			make([]decimal.Decimal, tradingPeriod),
			prevBid:				decimal.Zero(),
			overSold: 			rsiLowerLim,
			overBought: 		rsiUpperLim,
			lastOrderId:		"",
			readyToBuy:			true,
	}

	test(&bot)*/

	candle := candleBot{
		tradingPeriod: 	180, //180 minute candlestick
		tradesMade: 		0,
		numOfDecisions: 0,
		queue: 					[]candlestick{},
		readyToBuy: 		true,
    buyPrice:       decimal.Zero()
	}

	candle.fillQueue(3)

	for {
		candle.trade()
	}
}
