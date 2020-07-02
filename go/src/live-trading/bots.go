package main

import (
	"fmt"
	"time"
	"github.com/luno/luno-go/decimal"
)

/*TODO: Implement selling. Using stop loss order if there is a way to do it.
				Otherwise, implement this manually*/

// struct for the rsiBot
type rsiBot struct {
	tradingPeriod		int64										// How often the bot calculates a result
	tradesMade	   	int64						 				// total number of trades executed
	numOfDecisions 	int64					 					// number of times the bot calculates
	stopLoss		 	 	decimal.Decimal         // variable stop loss
	stopLossMult   	decimal.Decimal         // multiplier for stop loss
	pastBids			 	[]decimal.Decimal			  // cache of historical data from previous trading period
	overSold			 	decimal.Decimal					// bound to tell if the item is over sold
	overBought		 	decimal.Decimal					// bound to tell if the item is over bought
}

// function to execute buying of items
func buy(b *rsiBot) {
	currFunds := getAsset("GBP")
	price := 0.99 * getCurrAsk()
	buyableStock := currFunds.Div(price, 64) //WHY WAS THIS 8
	// checking if there are enough funds to buy the given amount of stock
	if currFunds.Sign() == 0 {
		fmt.Println("No funds available")
		return
	} else {
		//Create limit order

		req, err := PostLimitOrderRequest{
			Pair: pair,
			Price: price,
			Type: "BID", //We are putting in a bid to buy at the ask price
			Volume: buyableStock,
			//BaseAccountId: --> Not needed until using multiple strategies
			//CounterAccoundId: --> Same as above
			PostOnly: true
		}

		if err != nil {panic(err)}

		b.tradesMade++
		b.stopLoss = price
		// fmt.Println("Current funds: ",b.funds,"\n")
	}
}

// function to execute trades using the RSI bot
func (b *rsiBot) trade() {
	populatePastBids() //This also causes the program to sleep for one trading period

	// calculating RSI usig RSI algorithm
	rsi := rsi(pastBids)

	if getAsset("XBT").Sign() == 0 {
		if rsi.Cmp(b.overSold) == -1 {
			// buying stock if we have no stock and rsi is less than overSold bound
			buy(b)
			// printPortFolio(b.pf)
		}
		//Else, do nothing
	} else {
		bound := getCurrBid().Mul(b.stopLossMult)
		if bound.Cmp(b.stopLoss) == 1 {
			b.tradesMade++
			b.stopLoss = bound

			//TODO: Create a stop loss order

			// printPortFolio(b.pf)
		}
	}	//Else, do nothing

	b.numOfDecisions++
}

func printPortFolio(pf *portfolio) {
	fmt.Println("trade # :   ", pf.tradesMade)
	fmt.Println("funds : 			£", pf.funds)
	fmt.Println("stock : 			£", pf.stock)
	fmt.Println(".")
}
