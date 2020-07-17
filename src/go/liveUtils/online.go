package liveUtils

import (
	"context"
	"fmt"
	luno "github.com/luno/luno-go"
	"github.com/luno/luno-go/decimal"
	"time"
	. "TradingHackathon/src/go/rsi"
)

// function to execute buying of items
func buy(b *RsiBot, currAsk decimal.Decimal) {
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
		b.ReadyToBuy = false
		b.TradesMade++
		b.StopLoss = price
		b.BuyPrice = price
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

func sell(b *RsiBot, currBid decimal.Decimal) {
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
	b.ReadyToBuy = true
	b.TradesMade++
	for {
		time.Sleep(time.Minute)
		fmt.Println("Waiting for sell order to be partially filled")
		if funds.Cmp(getAsset("XBT")) == -1 {
			return
		}
	}
}

// function to execute trades using the RSI bot
func TradeLive(b *RsiBot) {

	time.Sleep(time.Minute)
	currAsk, currBid := getTicker()

	// calculating RSI using RSI algorithm
	var rsi decimal.Decimal
	rsi, b.UpEma, b.DownEma = GetRsi(b.PrevAsk, currAsk, b.UpEma, b.DownEma, b.TradingPeriod)
	fmt.Println("RSI", rsi, "U:", b.UpEma, "D:", b.DownEma)
	b.PrevAsk = currAsk

	if b.ReadyToBuy { // check if sell order has gone trough
		fmt.Println("Current Ask", currAsk)
		if rsi.Cmp(b.OverSold) == -1 && rsi.Sign() != 0 {
			buy(b, currAsk)
		}
	} else {
		bound := currBid.Mul(b.StopLossMult)

		fmt.Println("Current Bid", currBid)
		fmt.Println("Stop Loss", b.StopLoss)

		if (currBid.Cmp(b.BuyPrice) == 1 && currBid.Cmp(b.StopLoss) == -1) ||
			currBid.Cmp(b.BuyPrice.Mul(decimal.NewFromFloat64(0.98, 8))) == -1 {
			sell(b, currBid)
		} else if bound.Cmp(b.StopLoss) == 1 {
			b.StopLoss = bound
			fmt.Println("Stoploss changed to: ", b.StopLoss)
		}

	}
	b.NumOfDecisions++

}

func printPortFolio(b *RsiBot) {
	fmt.Println("trade # :   ", b.TradesMade)
	fmt.Println("funds : 			£", getAsset("GBP"))
	fmt.Println("stock : 		BTC", getAsset("XBT"))
}
