package main

import (
	"context"
	"fmt"
	luno "github.com/luno/luno-go"
	"github.com/luno/luno-go/decimal"
	"time"
)

// function to cancel most recent order
func cancelPrevOrder (b *offsetBot) {
	if b.PrevOrder == "" {return}
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
			fmt.Println("Failed to cancel previous order")
			cancelPrevOrder(b)
		}
	}
	fmt.Println("Previous order was filled. No need to cancel.")
}

// function to execute buying of items
func buy(b *offsetBot, currAsk decimal.Decimal) {
	cancelPrevOrder(b)
	time.Sleep(time.Second * 2)
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
		b.PrevOrder = res.OrderId
		b.readyToBuy = false
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

func sell(b *offsetBot, currBid decimal.Decimal) {
	cancelPrevOrder(b)
	time.Sleep(time.Second * 2)
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
	b.PrevOrder = res.OrderId
	b.readyToBuy = true
	for {
		time.Sleep(time.Minute)
		fmt.Println("Waiting for sell order to be partially filled")
		if funds.Cmp(getAsset("XBT")) == -1 {
			return
		}
	}
}
