package liveUtils

import (
	"context"
	"fmt"
	luno "github.com/luno/luno-go"
	"github.com/luno/luno-go/decimal"
	"time"
)

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
			Pair:   Pair,
			Price:  price,
			Type:   "BID", //We are putting in a bid to buy at the ask price
			Volume: buyableStock,
			//BaseAccountId: --> Not needed until using multiple strategies
			//CounterAccountId: --> Same as above
			PostOnly: true,
		}
		res, err := Client.PostLimitOrder(context.Background(), &req)
		for err != nil {
			fmt.Println(err)
			time.Sleep(time.Second * 30)
			res, err = Client.PostLimitOrder(context.Background(), &req)
		}
		fmt.Println("BUY - order ", res.OrderId, " placed at ", price)
		b.readyToBuy = false
		b.tradesMade++
		b.stopLoss = price
		b.buyPrice = price
		// wait till order has gone through
		for {
			time.Sleep(time.Minute)
			fmt.Println("Waiting for buy order to be partially filled")
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
		Pair:   Pair,
		Price:  price,
		Type:   "ASK", //We are putting in a ask to sell at the bid price
		Volume: volumeToSell,
		//BaseAccountId: --> Not needed until using multiple strategies
		//CounterAccoundId: --> Same as above
		PostOnly: true,
	}
	res, err := Client.PostLimitOrder(context.Background(), &req)
	for err != nil {
		fmt.Println(err)
		time.Sleep(time.Minute)
		res, err = Client.PostLimitOrder(context.Background(), &req)
	}

	fmt.Println("SELL - order ", res.OrderId, " placed at ", price)
	b.readyToBuy = true
	b.tradesMade++
	for {
		time.Sleep(time.Minute)
		fmt.Println("Waiting for sell order to be partially filled")
		if funds.Cmp(getAsset("XBT")) == -1 {
			return
		}
	}
}

// function to execute trades using the RSI bot
func tradeLive(b *rsiBot) {

	time.Sleep(time.Minute)
	currAsk, currBid := getTicker()

	// calculating RSI using RSI algorithm
	var rsi decimal.Decimal
	rsi, b.upEma, b.downEma = getRsi(b.prevAsk, currAsk, b.upEma, b.downEma, b.tradingPeriod)
	fmt.Println("RSI", rsi, "U:", b.upEma, "D:", b.downEma)
	b.prevAsk = currAsk

	if b.readyToBuy { // check if sell order has gone trough
		fmt.Println("Current Ask", currAsk)
		if rsi.Cmp(b.overSold) == -1 && rsi.Sign() != 0 {
			buy(b, currAsk)
		}
	} else {
		bound := currBid.Mul(b.stopLossMult)

		fmt.Println("Current Bid", currBid)
		fmt.Println("Stop Loss", b.stopLoss)

		if (currBid.Cmp(b.buyPrice) == 1 && currBid.Cmp(b.stopLoss) == -1) ||
			currBid.Cmp(b.buyPrice.Mul(decimal.NewFromFloat64(0.98, 8))) == -1 {
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
