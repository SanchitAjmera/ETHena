package main

import (
	"context"
	"fmt"
	luno "github.com/luno/luno-go"
	"github.com/luno/luno-go/decimal"
	"time"
)

// struct for the rsiBot
type rsiBot struct {
	tradingPeriod  int64             // No of past asks used to calculate RSI
	tradesMade     int64             // total number of trades executed
	numOfDecisions int64             // number of times the bot calculates
	stopLoss       decimal.Decimal   // variable stop loss
	stopLossMult   decimal.Decimal   // multiplier for stop loss
	pastAsks       []decimal.Decimal // cache of historical data from previous trading period
	overSold       decimal.Decimal   // bound to tell the bot when to buy
	readyToBuy     bool              // false means ready to sell
	buyPrice       decimal.Decimal   // stores most recent price we bought at
}

// function to execute buying of items
func buy(b *rsiBot, currAsk decimal.Decimal) {
	time.Sleep(time.Second * 5)
	targetFunds, currFunds := getAssets("XRP", "XBT")
	price := currAsk.Sub(decimal.NewFromFloat64(0.00000001, 8))
	buyableStock := currFunds.Div(price, 8)
	buyableStock = buyableStock.ToScale(0)
	// checking if there are no funds available
	if currFunds.Sign() == 0 {
		fmt.Println("No funds available")
		return
	} else {
		//Create limit order
		req := luno.PostLimitOrderRequest{
			Pair:   pair,
			Price:  price,
			Type:   "BID", //We are putting in a bid to buy at the ask price
			Volume: buyableStock,
			//BaseAccountId: --> Not needed until using multiple strategies
			//CounterAccountId: --> Same as above
			PostOnly: true,
		}
		res, err := client.PostLimitOrder(context.Background(), &req)
		for err != nil {
			fmt.Println(err)
			time.Sleep(time.Second * 30)
			res, err = client.PostLimitOrder(context.Background(), &req)
		}
		fmt.Println("BUY - order ", res.OrderId, " placed at ", price)
		b.readyToBuy = false
		b.tradesMade++
		b.stopLoss = price
		b.buyPrice = price
		// wait till order has gone through
		for {
			time.Sleep(time.Minute)

			if targetFunds.Cmp(getAsset("XRP")) == -1 {
				return
			}
		}
	}
}

func sell(b *rsiBot, currBid decimal.Decimal) {
	time.Sleep(time.Second * 5)
	volumeToSell, funds := getAssets("XRP","XBT")
	price := currBid.Add(decimal.NewFromFloat64(0.00000001, 8))
	req := luno.PostLimitOrderRequest{
		Pair:   pair,
		Price:  price,
		Type:   "ASK", //We are putting in a ask to sell at the bid price
		Volume: volumeToSell,
		//BaseAccountId: --> Not needed until using multiple strategies
		//CounterAccoundId: --> Same as above
		PostOnly: true,
	}
	res, err := client.PostLimitOrder(context.Background(), &req)
	for err != nil {
		fmt.Println(err)
		time.Sleep(time.Minute)
		res, err = client.PostLimitOrder(context.Background(), &req)
	}

	fmt.Println("SELL - order ", res.OrderId, " placed at ", price)
	b.readyToBuy = true
	b.tradesMade++
	for {
		time.Sleep(time.Minute)
		if funds.Cmp(getAsset("XBT")) == -1 {
			return
		}
	}
}

// function to execute trades using the RSI bot
func (b *rsiBot) trade() {

	time.Sleep(time.Minute)
	currAsk, currBid := getTicker()
	b.pastAsks = append(b.pastAsks[1:], currAsk)
	// calculating RSI using RSI algorithm
	rsi := rsi(b.pastAsks)

	if b.readyToBuy { // check if sell order has gone trough
		fmt.Println("Current Ask", currAsk)
		fmt.Println("RSI", rsi)
		if rsi.Cmp(b.overSold) == -1 && rsi.Sign() != 0 {
			buy(b, currAsk)
		}
	} else {
		bound := currBid.Mul(b.stopLossMult)

		fmt.Println("Current Bid", currBid)
		fmt.Println("Stop Loss", b.stopLoss)

		if (currBid.Cmp(b.buyPrice) == 1 && currBid.Cmp(b.stopLoss) == -1) ||
			currBid.Cmp(b.buyPrice.Mul(decimal.NewFromFloat64(0.95, 8))) == -1 {
			sell(b, currBid)
		} else if bound.Cmp(b.stopLoss) == 1 {
			b.stopLoss = bound
			fmt.Println("Stoploss changed to: ", b.stopLoss)
		}

	}
	b.numOfDecisions++

}

func printPortFolio(b *rsiBot) {
	fmt.Println("trade # :   ", b.tradesMade)
	fmt.Println("funds : 			£", getAsset("GBP"))
	fmt.Println("stock : 		BTC", getAsset("XBT"))
}
