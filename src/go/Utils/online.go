package Utils

import (
	"context"
	luno "github.com/luno/luno-go"
	"github.com/luno/luno-go/decimal"
//	"log"
	time "time"
)
// global variables for printing purposes
var stopLossValues []decimal.Decimal
var macdValues []decimal.Decimal
var rsiValues []decimal.Decimal
var scoreValues []decimal.Decimal

// function to cancel most recent order
func cancelPrevOrder(b *RsiBot) {
	PrintStatus(b, decimal.Zero(), decimal.Zero(), "CANCELLING PREVIOUS ORDER","", []([]decimal.Decimal){rsiValues, macdValues, stopLossValues, scoreValues})
	if b.PrevOrder == "" {
	//	log.Println("No previous order to cancel")
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
		//	log.Println("Successfully cancelled previous order")
			PrintStatus(b, decimal.Zero(), decimal.Zero(), "CANCELLED PREVIOUS ORDER","", []([]decimal.Decimal){rsiValues, macdValues, stopLossValues, scoreValues})

		} else {
		//	log.Println("ERROR! Failed to cancel previous order")
			cancelPrevOrder(b)
		}
	}
	PrintStatus(b, decimal.Zero(), decimal.Zero(), "CANCELLED PREVIOUS ORDER","", []([]decimal.Decimal){rsiValues, macdValues, stopLossValues, scoreValues})
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
			PrintStatus(b, decimal.Zero(), currAsk, "NO FUND AVAILABLE","", []([]decimal.Decimal){rsiValues, macdValues, stopLossValues, scoreValues})
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
		//log.Println(err)
		time.Sleep(time.Second * 30)
		res, err = Client.PostLimitOrder(context.Background(), &req)
	}
	info := "BUY ORDER " + res.OrderId+  " PLACED AT " + price.String()[:6]
	PrintStatus(b, decimal.Zero(), currAsk, "BUY ORDER PLACED",info, []([]decimal.Decimal){rsiValues, macdValues, stopLossValues, scoreValues})

	b.PrevOrder = res.OrderId
	b.ReadyToBuy = false
	b.TradesMade++
	b.StopLoss = price
	b.BuyPrice = price
	// wait till order has gone through
	PrintStatus(b, decimal.Zero(), currAsk, "WAITING FOR BUY ORDER TO BE FILLED",info, []([]decimal.Decimal){rsiValues, macdValues, stopLossValues, scoreValues})
	counter := 0
	for {
		time.Sleep(2 * time.Second)
		counter++
		if startStock.Cmp(getAsset("ETH")) == -1 {
			PrintStatus(b, decimal.Zero(), currAsk, "BUY ORDER FILLED",info, []([]decimal.Decimal){rsiValues, macdValues, stopLossValues, scoreValues})
			return
		}
		if counter > 15 {
			b.TradesMade--
			PrintStatus(b, decimal.Zero(), currAsk, "RETRYING BUY ORDER",info, []([]decimal.Decimal){rsiValues, macdValues, stopLossValues, scoreValues})
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
	//log.Println("startstock after scale: ", startStock)
	if startStock.Sign() == 0 {
		PrintStatus(b, decimal.Zero(), decimal.Zero(), "NO FUND AVAILABLE","", []([]decimal.Decimal){rsiValues, macdValues, stopLossValues, scoreValues})
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
		//log.Println(err)
		time.Sleep(2 * time.Second)
		res, err = Client.PostLimitOrder(context.Background(), &req)
	}

	info := "SELL ORDER " + res.OrderId +  " PLACED AT " + price.String()[:6]
	PrintStatus(b, price, decimal.Zero(), "SELL ORDER PLACED",info, []([]decimal.Decimal){rsiValues, macdValues, stopLossValues, scoreValues})
	b.PrevOrder = res.OrderId
	b.ReadyToBuy = true
	b.TradesMade++
	b.BuyPrice = decimal.Zero()
	PrintStatus(b, price, decimal.Zero(), "WAITING FOR SELL ORDER TO BE FILLED",info, []([]decimal.Decimal){rsiValues, macdValues, stopLossValues, scoreValues})
	counter := 0
	for {
		time.Sleep(2 * time.Second)
		counter++
		if startFunds.Cmp(getAsset("XBT")) == -1 {
			PrintStatus(b, price, decimal.Zero(), "SELL ORDER FILLED",info, []([]decimal.Decimal){rsiValues, macdValues, stopLossValues, scoreValues})
			return
		}
		if counter > 15 {
			b.TradesMade--
			PrintStatus(b, price, decimal.Zero(), "RETRYING SELL ORDER",info, []([]decimal.Decimal){rsiValues, macdValues, stopLossValues, scoreValues})
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
	//printing
	var status string
	if (b.ReadyToBuy){
		status = "READY TO BUY"
	} else {
		status = "READY TO SELL"
	}

	b.Stack = append(b.Stack, stick)
	b.Stack = b.Stack[1:]
	prevema := Sma(b.PastAsks[b.LongestTradingPeriod-b.OffsetTraingPeriod : b.LongestTradingPeriod-1])

	scores := []decimal.Decimal{}

	var rsi decimal.Decimal
	rsi, b.UpEma, b.DownEma = GetRsi(b.PrevAsk, currAsk, b.UpEma, b.DownEma, b.RSITradingPeriod)
	rsiValues = append(rsiValues, rsi)

	rsiweighting := int(b.BotString[0])
	MACDweighting := int(b.BotString[1])
	Candlestickweighting := int(b.BotString[2])
	Offsetweighting := int(b.BotString[3])

	if rsiweighting != '0' && rsi != decimal.Zero(){
		rsiScore := decimal.NewFromInt64(100).Sub(rsi)
		for i := 0; i < rsiweighting; i++ {
			scores = append(scores, rsiScore)
		}
	}


	if MACDweighting != '0' {
		b.MACDlongperiodavg = Sma(b.PastAsks[b.LongestTradingPeriod-b.MACDTradingPeriodLR:])
		b.MACDshortperiodavg = Sma(b.PastAsks[b.LongestTradingPeriod-b.MACDTradingPeriodSR:])
		currdifference := b.MACDshortperiodavg.Sub(b.MACDlongperiodavg)
		macdValues = append(macdValues, currdifference)
		macdScore := decimal.NewFromInt64(100).Sub(currdifference.Div(decimal.NewFromFloat64(0.000001, 16), 16))
		for i := 0; i < MACDweighting; i++ {
			scores = append(scores, macdScore)
		}
	} else {
			macdValues = append(macdValues, decimal.Zero())
	}

	if Candlestickweighting != '0' {
		if Rev123(b.Stack[b.LongestTradingPeriod-3], b.Stack[b.LongestTradingPeriod-2], b.Stack[b.LongestTradingPeriod-1]) || Hammer(b.Stack[b.LongestTradingPeriod-1]) || InverseHammer(b.Stack[b.LongestTradingPeriod-1]) || WhiteSlaves(b.Stack[b.LongestTradingPeriod-3], b.Stack[b.LongestTradingPeriod-2], b.Stack[b.LongestTradingPeriod-1]) || MorningStar(b.Stack[b.LongestTradingPeriod-3], b.Stack[b.LongestTradingPeriod-2], b.Stack[b.LongestTradingPeriod-1]) {
			candlestickscore := decimal.NewFromInt64(100)
			for i := 0; i < Candlestickweighting; i++ {
				scores = append(scores, candlestickscore)
			}
		}
	}
	if Offsetweighting != '0' {
		ema := Ema(prevema, currAsk, b.OffsetTraingPeriod)
		if currAsk.Cmp(ema.Sub(b.Offset)) == -1 {
			offsetscore := decimal.NewFromInt64(100)
			for i := 0; i < Offsetweighting; i++ {
				scores = append(scores, offsetscore)
			}
		}
	}

	b.PastAsks = b.PastAsks[1:]
	b.PastAsks = append(b.PastAsks, currAsk)
	b.PrevAsk = currAsk

	averageScore := Sma(scores)

	scoreValues = append(scoreValues, averageScore)
	stopLossValues = append(stopLossValues, b.StopLoss)

	if len(rsiValues) > 5 {
		rsiValues = rsiValues[1:]
		stopLossValues = stopLossValues[1:]
		macdValues = macdValues[1:]
		scoreValues = scoreValues[1:]
	}
	PrintStatus(b, currBid, currAsk, status,"", []([]decimal.Decimal){rsiValues, macdValues, stopLossValues, scoreValues})

//	fmt.Println("Average Score: ", averageScore)
	PopulateFile(b, currAsk, currBid, rsi)

	if b.ReadyToBuy { // check if sell order has gone trough
		//log.Println("Current Ask", currAsk)
		//log.Println("RSI: ", rsi, "U: ", b.UpEma, "D: ", b.DownEma)
		if averageScore.Cmp(decimal.NewFromInt64(80)) == 1 {
			buy(b, currAsk)
		}
	} else {
		bound := currBid.Mul(b.StopLossMult)

	//	log.Println("Current Bid", currBid)
		//log.Println("Stop Loss", b.StopLoss)

		if (currBid.Cmp(b.BuyPrice) == 1 && currBid.Cmp(b.StopLoss) == -1) ||
			currBid.Cmp(b.BuyPrice.Mul(decimal.NewFromFloat64(0.99, 8))) == -1 {
			sell(b, currBid)
		} else if bound.Cmp(b.StopLoss) == 1 {
			b.StopLoss = bound
	//		log.Println("Stoploss changed to: ", b.StopLoss)
		}
		b.NumOfDecisions++
	}
}
