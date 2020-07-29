package liveUtils

import (
	. "TradingHackathon/src/go/rsi"
	"context"
	"fmt"
	"time"

	luno "github.com/luno/luno-go"
	"github.com/luno/luno-go/decimal"
)

// function to cancel most recent order
func cancelPrevOrder(b *RsiBot) {
	if b.PrevOrder == "" {
		fmt.Println("No previous order to cancel")
		return
	}
	time.Sleep(time.Second * 2)
	checkReq := luno.GetOrderRequest{Id: b.PrevOrder}
	checkRes, err := Client.GetOrder(context.Background(), &checkReq)
	if err != nil {
		panic(err)
	}
	if checkRes.State == "PENDING" {
		time.Sleep(time.Second * 2)
		req := luno.StopOrderRequest{OrderId: b.PrevOrder}
		res, err := Client.StopOrder(context.Background(), &req)
		if err != nil {
			panic(err)
		}
		if res.Success {
			fmt.Println("Successfully cancelled previous order")
		} else {
			fmt.Println("ERROR! Failed to cancel previous order")
			cancelPrevOrder(b)
		}
	}
	fmt.Println("Previous order was filled. Cancellation not required.")
}

// function to execute buying of items
func buy(b *RsiBot, currAsk decimal.Decimal) {
	cancelPrevOrder(b)
	time.Sleep(time.Second * 2)
	startStock, startFunds := getAssets("ETH", "XBT")
	fmt.Println("startFunds: ", startFunds)
	fmt.Println("StartStock: ", startStock)
	fmt.Println("currAsk: ", currAsk)
	price := currAsk.Sub(decimal.NewFromFloat64(0.000001, 8))
	fmt.Println("price: ", price)
	buyableStock := startFunds.Div(price, 8)
	fmt.Println("buyablestock before scale: ", buyableStock)
	buyableStock = buyableStock.ToScale(2)
	// checking if there are no funds available
	fmt.Println("buyablestock after scale: ", buyableStock)
	if buyableStock.Sign() == 0 {
		fmt.Println("Not enough funds available")
		return
	}
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
	b.PrevOrder = res.OrderId
	b.ReadyToBuy = false
	b.TradesMade++
	b.StopLoss = price
	b.BuyPrice = price
	// wait till order has gone through
	fmt.Println("Waiting for buy order to be partially filled")
	for {
		time.Sleep(2 * time.Second)
		if startStock.Cmp(getAsset("ETH")) == -1 {
			fmt.Println("Buy order has been partially filled")
			return
		}
	}
}

func sell(b *RsiBot, currBid decimal.Decimal) {
	cancelPrevOrder(b)
	time.Sleep(time.Second * 2)
	startStock, startFunds := getAssets("ETH", "XBT")
	price := currBid.Add(decimal.NewFromFloat64(0.000001, 8))
	req := luno.PostLimitOrderRequest{
		Pair:   Pair,
		Price:  price,
		Type:   "ASK", //We are putting in a ask to sell at the bid price
		Volume: startStock,
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
	b.PrevOrder = res.OrderId
	b.ReadyToBuy = true
	b.TradesMade++
	fmt.Println("Waiting for sell order to be partially filled")
	for {
		time.Sleep(2 * time.Second)
		if startFunds.Cmp(getAsset("XBT")) == -1 {
			fmt.Println("Sell order has been partially filled")
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
	b.PrevAsk = currAsk

	PopulateFile(b, currAsk, currBid, rsi)

	if b.ReadyToBuy { // check if sell order has gone trough
		fmt.Println("Current Ask", currAsk)
		fmt.Println("RSI: ", rsi, "U: ", b.UpEma, "D: ", b.DownEma)
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
