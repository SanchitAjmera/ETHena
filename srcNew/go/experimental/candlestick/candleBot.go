package main

import (
  "fmt"
  "github.com/luno/luno-go/decimal"
)



type candleBot struct {
	TradingPeriod		int64										// How often the bot calculates a result
	TradesMade	   	int64						 				// total number of trades executed
	NumOfDecisions 	int64					 					// number of times the bot calculates
	queue					  []candlestick				  	// previous 3 candlesticks
  ReadyToBuy      bool                    // holds the state of the bot
	BuyPrice				decimal.Decimal					// price we bought at
  currRow         int64
}

type candlestick struct {
	openAsk   decimal.Decimal
	closeAsk  decimal.Decimal
	maxAsk    decimal.Decimal
	minAsk    decimal.Decimal
	openBid   decimal.Decimal
	closeBid  decimal.Decimal
	maxBid    decimal.Decimal
	minBid    decimal.Decimal
}

func (b *candleBot) getCurrCandlestick() candlestick{
  var maxInt64 int64 = 1844674407370955200
  result := candlestick{
    openAsk:  decimal.Zero(),
    closeAsk: decimal.Zero(),
    maxAsk:   decimal.Zero(),
    minAsk:   decimal.NewFromInt64(maxInt64),
    openBid:  decimal.Zero(),
    closeBid:  decimal.Zero(),
    maxBid:   decimal.Zero(),
    minBid:   decimal.NewFromInt64(maxInt64),
  }

  result.openBid = getBid(b.currRow)
  result.openAsk = getAsk(b.currRow)
  b.currRow += 1

	callsPerMinute := 1

  for i := 1; i < callsPerMinute * int(b.TradingPeriod)-1; i++ {

    currBid := getBid(b.currRow)
    currAsk := getAsk(b.currRow)

    if result.maxAsk.Cmp(currAsk) == -1 {
      result.maxAsk = currAsk
    }

    if result.maxBid.Cmp(currBid) == -1 {
      result.maxBid = currBid
    }

    if currAsk.Cmp(result.minAsk) == -1 {
      result.minAsk = currAsk
    }

    if currBid.Cmp(result.minBid) == -1 {
      result.minBid = currBid
    }

    b.currRow += 1
  }

  result.closeBid = getBid(b.currRow+int64(b.TradingPeriod))
  result.closeAsk = getAsk(b.currRow+int64(b.TradingPeriod))

  b.currRow += 1

  // fmt.Println("Bids:  Open - ",result.openBid," High - ",result.maxBid," Low - ",result.minBid," Close - ",result.closeBid)
  // fmt.Println("Asks:  Open - ",result.openAsk," High - ",result.maxAsk," Low - ",result.minAsk," Close - ",result.closeAsk)

  return result
}

func (b *candleBot) fillQueue(queueSize int) {
  for i := 0; i < queueSize; i++ {
    fmt.Println("Filling queue: ",i+1,"/",queueSize,"\n")
    b.queue = append(b.queue, b.getCurrCandlestick())
  }
}

func (b *candleBot) trade3() {
		b1Op := b.queue[0].openAsk
		b1Cl := b.queue[0].closeAsk
		b1Max := b.queue[0].maxAsk
		b1Min := b.queue[0].minAsk

		b2Op := b.queue[1].openAsk
		b2Cl := b.queue[1].closeAsk
		b2Max := b.queue[1].maxAsk
		b2Min := b.queue[1].minAsk

		b3Op := b.queue[2].openAsk
		b3Cl := b.queue[2].closeAsk
		b3Max := b.queue[2].maxAsk
		b3Min := b.queue[2].minAsk

		if b2Max.Cmp(b1Max) == 1 && b2Max.Cmp(b3Max) == 1 && b2Min.Cmp(b1Min) == 1 && b2Min.Cmp(b3Min) == 1 {
			if b1Cl.Cmp(b1Op) == 1 && b3Op.Cmp(b3Cl) == 1 && b2Cl.Cmp(b2Op) == 1 {
        b.sell()
			}
		} else if b2Max.Cmp(b1Max) == -1 && b2Max.Cmp(b3Max) == -1 && b2Min.Cmp(b1Min) == -1 && b2Min.Cmp(b3Min) == -1{
			if b1Cl.Cmp(b1Op) == -1 && b3Op.Cmp(b3Cl) == -1 && b2Cl.Cmp(b2Op) == -1{
        b.buy()
			}
		} else {
  //     fmt.Println("HOLD at",b3Cl)
    }
    b.NumOfDecisions++
}

func (b *candleBot) trade() {
//  fmt.Println("currRow", b.currRow)
  //Move the queue forward
  for i := 0; i < len(b.queue) - 1; i++ {
    b.queue[i] = b.queue[i+1]
  }
  b.queue[len(b.queue) - 1] = b.getCurrCandlestick()

  b.trade3()
}


func (b *candleBot) buy() {
  price := getAsk(b.currRow).Mul(decimal.NewFromFloat64(0.99999,8))

//  buyableStock := currFunds.Div(price, 8)
  // checking if there are enough funds to buy the given amount of stock
/*  if currFunds.Sign() == 0 {
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
    if err != nil {panic(err)}*/
    fmt.Println("BUYS at", price, " currRow:", b.currRow)
    b.ReadyToBuy = false
		b.BuyPrice = price
    b.TradesMade++

}


func (b *candleBot) sell() {
	price := getBid(b.currRow).Mul(decimal.NewFromFloat64(1.00001, 8))
/*
	if price < b.BuyPrice {
		fmt.Println("Spread too high to sell")
		return
	}

	req := luno.PostLimitOrderRequest{
		Pair: pair,
		Price: price,
		Type: "ASK", //We are putting in a ask to sell at the bid price
		Volume: getAsset("XBT"),
		//BaseAccountId: --> Not needed until using multiple strategies
		//CounterAccoundId: --> Same as above
		PostOnly: true,
	}

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

	res, err := client.PostLimitOrder(context.Background(), &req)
	if err != nil {panic(err)}
*/
	fmt.Println("SELL at", price)
//  fmt.Println("At the time above, ask price was:",currAsk,"\n")
	//b.lastOrderId = res.OrderId
	b.ReadyToBuy = true
	b.TradesMade++
}
