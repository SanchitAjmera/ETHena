package liveUtils

import (
	. "TradingHackathon/src/go/rsi"
	"context"
	"log"
	"time"

	luno "github.com/luno/luno-go"
	"github.com/luno/luno-go/decimal"
)

// function to cancel most recent order
func cancelPrevOrder(b *RsiBot) {
	if b.PrevOrder == "" {
		log.Println("No previous order to cancel")
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
			log.Println("Successfully cancelled previous order")
			return
		} else {
			log.Println("ERROR! Failed to cancel previous order")
			cancelPrevOrder(b)
		}
	}
	log.Println("Previous order was filled. Cancellation not required.")
}

// function to execute buying of items
func buy(b *RsiBot, currAsk decimal.Decimal) {
	cancelPrevOrder(b)
	time.Sleep(time.Second * 2)

	startStock, startFunds := getAssets("ETH", "XBT")
	log.Println("startFunds: ", startFunds)
	log.Println("StartStock: ", startStock)
	log.Println("currAsk: ", currAsk)
	price := currAsk.Sub(decimal.NewFromFloat64(0.000001, 8))
	log.Println("price: ", price)
	buyableStock := startFunds.Div(price, 8)
	log.Println("buyablestock before scale: ", buyableStock)
	buyableStock = buyableStock.ToScale(2)

	// checking if there are no funds available
	log.Println("buyablestock after scale: ", buyableStock)
	if buyableStock.Sign() == 0 {
		log.Println("Not enough funds available")
		b.ReadyToBuy = false
		return
	}
	//Create limit order
	req := luno.PostLimitOrderRequest{
		Pair:   PairName,
		Price:  price,
		Type:   "BID", //We are putting in a bid to buy at the ask price
		Volume: buyableStock,
		//BaseAccountId: --> Not needed until using multiple strategies
		//CounterAccountId: --> Same as above
		PostOnly: true,
	}
	res, err := Client.PostLimitOrder(context.Background(), &req)
	for err != nil {
		log.Println(err)
		time.Sleep(time.Second * 30)
		res, err = Client.PostLimitOrder(context.Background(), &req)
	}
	log.Println("BUY - order ", res.OrderId, " placed at ", price)
	b.PrevOrder = res.OrderId
	b.ReadyToBuy = false
	b.TradesMade++
	b.BuyPrice = price
	// wait till order has gone through
	log.Println("Waiting for buy order to be partially filled")
	counter := 0
	for {
		time.Sleep(2 * time.Second)
		if startStock.Cmp(getAsset("ETH")) == -1 {
			log.Println("Buy order has been partially filled")
			return
		}
		counter++
		if counter > 15 {
			log.Println("Timeout. Retrying buy")
			time.Sleep(2 * time.Second)
			b.TradesMade--
			buy(b, getTickerRes().Ask)
			return
		}
	}
}

func sell(b *RsiBot, currBid decimal.Decimal) {
	cancelPrevOrder(b)
	time.Sleep(time.Second * 2)

	startStock, startFunds := getAssets("ETH", "XBT")
	startStock = startStock.ToScale(2)
	price := currBid.Add(decimal.NewFromFloat64(0.000001, 8))

	req := luno.PostLimitOrderRequest{
		Pair:   PairName,
		Price:  price,
		Type:   "ASK", //We are putting in a ask to sell at the bid price
		Volume: startStock,
		//BaseAccountId: --> Not needed until using multiple strategies
		//CounterAccoundId: --> Same as above
		PostOnly: true,
	}
	res, err := Client.PostLimitOrder(context.Background(), &req)
	for err != nil {
		log.Println(err)
		time.Sleep(2 * time.Second)
		res, err = Client.PostLimitOrder(context.Background(), &req)
	}

	log.Println("SELL - order ", res.OrderId, " placed at ", price)
	b.PrevOrder = res.OrderId
	b.ReadyToBuy = true
	b.TradesMade++
	log.Println("Waiting for sell order to be partially filled")
	counter := 0
	for {
		time.Sleep(2 * time.Second)
		if startFunds.Cmp(getAsset("XBT")) == -1 {
			log.Println("Sell order has been partially filled")
			return
		}
		counter++
		if counter > 15 {
			log.Println("Timeout. Retrying sell")
			time.Sleep(2 * time.Second)
			b.TradesMade--
			sell(b, getTickerRes().Bid)
			return
		}
	}
}

// function to execute trades using the RSI bot
func TradeLive(b *RsiBot) {
	time.Sleep(5 * time.Second)
	res := getTickerRes()
	currAsk, currBid := res.Ask, res.Bid

	// calculating RSI using RSI algorithm
	var rsi decimal.Decimal
	rsi, b.UpEma, b.DownEma = GetRsi(b.PrevAsk, currAsk, b.UpEma, b.DownEma, b.TradingPeriod)
	b.PrevAsk = currAsk

	PopulateFile(b, currAsk, currBid, rsi)

	if b.ReadyToBuy { // check if sell order has gone trough
		log.Println("Current Ask", currAsk)
		log.Println("RSI: ", rsi, "U: ", b.UpEma, "D: ", b.DownEma)
		if rsi.Cmp(b.OverSold) == -1 && rsi.Sign() != 0 {
			buy(b, currAsk)
		}
	} else {


		log.Println("Current Bid", currBid)

		if (currBid.Cmp(b.BuyPrice) == 1) ||
			currBid.Cmp(b.BuyPrice.Mul(decimal.NewFromFloat64(0.9955, 8))) == -1 {
			sell(b, currBid)
		} 

	}
	b.NumOfDecisions++

}
