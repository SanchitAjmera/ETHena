package Utils

import (
	"context"
	"fmt"
	luno "github.com/luno/luno-go"
	"github.com/luno/luno-go/decimal"
	"log"
	"os"
	time "time"
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
	price := currAsk.Sub(decimal.NewFromFloat64(0.000001, 8))
	buyableStock := startFunds.Div(price, 8)
	buyableStock = buyableStock.ToScale(2)

	// checking if there are no funds available

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
	b.StopLoss = price
	b.BuyPrice = price
	// wait till order has gone through
	log.Println("Waiting for buy order to be partially filled")
	counter := 0
	for {
		time.Sleep(2 * time.Second)
		counter++
		if startStock.Cmp(getAsset("ETH")) == -1 {
			log.Println("Buy order has been partially filled")
			return
		}
		if counter > 15 {
			b.TradesMade--
			log.Println("Buy order timed out. Retrying")
			buy(b, currAsk)
			return
		}
	}
}

func sell(b *RsiBot, currBid decimal.Decimal) {
	cancelPrevOrder(b)
	time.Sleep(time.Second * 2)

	startStock, startFunds := getAssets("ETH", "XBT")
	startStock = startStock.ToScale(2)

	// checking if there are no stock available
	log.Println("startstock after scale: ", startStock)
	if startStock.Sign() == 0 {
		log.Println("Not enough stock available")
		b.ReadyToBuy = true
		return
	}
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
		counter++
		if startFunds.Cmp(getAsset("XBT")) == -1 {
			log.Println("Sell order has been partially filled")
			return
		}
		if counter > 15 {
			b.TradesMade--
			log.Println("Sell order timed out. Retrying")
			sell(b, currBid)
			return
		}
	}
}

// TradeLive function to execute trades using the RSI bot
func TradeLive(b *RsiBot) {
	time.Sleep(time.Second)
	stick := GetCandleStick(b.TimeInterval)
	currAsk, currBid := stick.CloseAsk, stick.CloseBid
	b.Stack = append(b.Stack, stick)
	b.Stack = b.Stack[1:]
	prevema := Sma(b.PastAsks[b.LongestTradingPeriod-b.OffsetTraingPeriod : b.LongestTradingPeriod-1])

	botstring := ""
	botstring = os.Args[2]
	if botstring == "0000" {
		fmt.Println("No Strategies Chosen. Bot has been stopped")

	}
	scores := []decimal.Decimal{}

	var rsi decimal.Decimal
	rsi, b.UpEma, b.DownEma = GetRsi(b.PrevAsk, currAsk, b.UpEma, b.DownEma, b.RSITradingPeriod)
	if []rune(botstring)[0] == '1' {
		rsiScore := decimal.NewFromInt64(100).Sub(rsi)
		scores = append(scores, rsiScore)
	}

	b.PastAsks = b.PastAsks[1:]
	b.PastAsks = append(b.PastAsks, currAsk)
	b.PrevAsk = currAsk

	if []rune(botstring)[1] == '1' {
		b.MACDlongperiodavg = Sma(b.PastAsks[b.LongestTradingPeriod-b.MACDTradingPeriodLR:])
		b.MACDshortperiodavg = Sma(b.PastAsks[b.LongestTradingPeriod-b.MACDTradingPeriodSR:])
		currdifference := b.MACDshortperiodavg.Sub(b.MACDlongperiodavg)
		macdScore := decimal.NewFromInt64(100).Sub(currdifference.Div(decimal.NewFromFloat64(0.000001, 16), 16))
		scores = append(scores, macdScore)
	}

	if []rune(botstring)[2] == '1' {
		if Rev123(b.Stack[b.LongestTradingPeriod-3], b.Stack[b.LongestTradingPeriod-2], b.Stack[b.LongestTradingPeriod-1]) || Hammer(b.Stack[b.LongestTradingPeriod-1]) || InverseHammer(b.Stack[b.LongestTradingPeriod-1]) || WhiteSlaves(b.Stack[b.LongestTradingPeriod-3], b.Stack[b.LongestTradingPeriod-2], b.Stack[b.LongestTradingPeriod-1]) || MorningStar(b.Stack[b.LongestTradingPeriod-3], b.Stack[b.LongestTradingPeriod-2], b.Stack[b.LongestTradingPeriod-1]) {
			candlestickscore := decimal.NewFromInt64(100)
			scores = append(scores, candlestickscore)
		}

	}
	if []rune(botstring)[3] == '1' {
		ema := ema(prevema, currAsk, b.OffsetTraingPeriod)
		if currAsk.Cmp(ema.Sub(b.Offset)) == -1 {
			offsetscore := decimal.NewFromInt64(100)
			scores = append(scores, offsetscore)
		}

	}
	averageScore := Sma(scores)
	fmt.Println("Average Score: ", averageScore)
	PopulateFile(b, currAsk, currBid, rsi)

	if b.ReadyToBuy { // check if sell order has gone trough
		//log.Println("Current Ask", currAsk)
		//log.Println("RSI: ", rsi, "U: ", b.UpEma, "D: ", b.DownEma)
		if averageScore.Cmp(decimal.NewFromInt64(80)) == 1 {
			buy(b, currAsk)
		}
	} else {
		bound := currBid.Mul(b.StopLossMult)

		log.Println("Current Bid", currBid)
		log.Println("Stop Loss", b.StopLoss)

		if (currBid.Cmp(b.BuyPrice) == 1 && currBid.Cmp(b.StopLoss) == -1) ||
			currBid.Cmp(b.BuyPrice.Mul(decimal.NewFromFloat64(0.99, 8))) == -1 {
			sell(b, currBid)
		} else if bound.Cmp(b.StopLoss) == 1 {
			b.StopLoss = bound
			log.Println("Stoploss changed to: ", b.StopLoss)
		}
		b.NumOfDecisions++
	}
}
