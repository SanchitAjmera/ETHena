package main

import (
	"fmt"
	"time"
	"context"
  luno "github.com/luno/luno-go"
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
	pastAsks			 	[]decimal.Decimal			  // cache of historical data from previous trading period
	prevBid					decimal.Decimal					// previous bid price
	overSold			 	decimal.Decimal					// bound to tell if the item is over sold
	overBought		 	decimal.Decimal					// bound to tell if the item is over bought
	lastOrderId			string									// most recent sell order id
	readyToBuy			bool										// false means ready to sell
}

// function to execute buying of items
func buy(b *rsiBot) {
	currFunds := getAsset("GBP")
	price := getCurrAsk().Mul(decimal.NewFromFloat64(0.99999,8))
	buyableStock := currFunds.Div(price, 8)
	// checking if there are enough funds to buy the given amount of stock
	if currFunds.Sign() == 0 {
		fmt.Println("No funds available")
		return
	} else {
		//Create limit order

		req := luno.PostLimitOrderRequest{
			Pair: pair,
			Price: price,
			Type: "BID", //We are putting in a bid to buy at the ask price
			Volume: buyableStock,
			//BaseAccountId: --> Not needed until using multiple strategies
			//CounterAccoundId: --> Same as above
			PostOnly: true,
		}
		res, err := client.PostLimitOrder(context.Background(), &req)
		if err != nil {panic(err)}
		fmt.Println("BUY - order ", res.OrderId, " placed at ", price)
		b.readyToBuy = false
		b.tradesMade++
		b.stopLoss = price
	}
}

func sell(b *rsiBot) {
	price := getCurrBid().Mul(decimal.NewFromFloat64(1.00001, 8))
	req := luno.PostLimitOrderRequest{
		Pair: pair,
		Price: price,
		Type: "ASK", //We are putting in a ask to sell at the bid price
		Volume: getAsset("XBT"),
		//BaseAccountId: --> Not needed until using multiple strategies
		//CounterAccoundId: --> Same as above
		PostOnly: true,
	}
	res, err := client.PostLimitOrder(context.Background(), &req)
	if err != nil {panic(err)}

	if (b.lastOrderId != "") {
		stopReq := luno.StopOrderRequest{OrderId: b.lastOrderId}
		stopRes, err := client.StopOrder(context.Background(), &stopReq)
		if err != nil {panic(err)}
		if !stopRes.Success {
			fmt.Println("Failed to cancel order")
			return
		}
		fmt.Println("Previous order successfully cancelled")
	}

	fmt.Println("SELL - order ", res.OrderId, " placed at ", price)
	b.lastOrderId = res.OrderId
	b.readyToBuy = true
	b.tradesMade++

}

func populatePastAsks (b *rsiBot) {
	//Populating past asks with 1 tradingPeriod worth of data
	var i int64 = 0
	for i < b.tradingPeriod {
		b.pastAsks[i] = getCurrAsk()

		buffer := ""
		if (i < 9) {buffer = " "}

		fmt.Println("Filling past asks: ",buffer, i+1,"/",b.tradingPeriod,":  £",b.pastAsks[i])
		i++

		time.Sleep(time.Minute) // Change to minute
	}
	fmt.Println("")
}

// function to execute trades using the RSI bot
func (b *rsiBot) trade(){
	// calculating RSI usig RSI algorithm
	rsi := rsi(b.pastAsks)

	if b.readyToBuy {
		if rsi.Cmp(b.overSold) == -1 {
			// buying stock if we have no stock and rsi is less than overSold bound
			//Delete from here till next comment
			currAsk := getCurrAsk()
			b.readyToBuy = false
			price := currAsk.Mul(decimal.NewFromFloat64(0.99999,8))
			b.stopLoss = price
			fmt.Println("BUYING at ", price)
			//buy(b)
			printPortFolio(b)
		}
		//Else, do nothing
	} else {

	 	currBid := getCurrBid()
		bound := currBid.Mul(b.stopLossMult)

		if b.prevBid.Cmp(b.stopLoss) == 1 && currBid.Cmp(b.stopLoss) == -1 {
			//Delete from here till next comment
			price := currBid.Mul(decimal.NewFromFloat64(1.00001, 8))
			fmt.Println("SELLING at ", price)
			b.readyToBuy = true

			//sell(b)
			printPortFolio(b)
		} else if bound.Cmp(b.stopLoss) == 1 {
			b.stopLoss = bound
			fmt.Println("Stoploss: ",b.stopLoss)
		} //Else, do nothing
	}
	b.numOfDecisions++

	time.Sleep(time.Minute)

	var returnVal decimal.Decimal
	b.prevBid , returnVal, _ = getTicker()
	b.pastAsks = append(b.pastAsks, returnVal)
	b.pastAsks = b.pastAsks[1:]
	//fmt.Println("Current Bid: £", returnVal)

}
