package backtestingUtils

import (
	"fmt"
	"github.com/luno/luno-go/decimal"
)

var xrp decimal.Decimal
var xbt decimal.Decimal
var currRow int64

func initialiseFunds(xbtFunds decimal.Decimal, xrpStock decimal.Decimal) {
	parseXlsx()
	xrp = xrpStock
	xbt = xbtFunds
	currRow = 0
}

// function to execute buying of items
func buyOffline(b *rsiBot, currAsk decimal.Decimal) {
	targetFunds, currFunds := xrp, xbt
	price := currAsk.Sub(decimal.NewFromFloat64(0.00000001, 8))
	buyableStock := currFunds.Div(price, 8)
	buyableStock = buyableStock.ToScale(0)
	// checking if there are no funds available
	if currFunds.Sign() == 0 {
		fmt.Println("No funds available")
		return
	} else {
		fmt.Println("BUY order placed at", price)
		b.readyToBuy = false
		b.tradesMade++
		b.stopLoss = price
		b.buyPrice = price

		// update funds
		xbt = xbt.Sub(price)
		xrp = xrp.Add(buyableStock)
	}
}

func sellOffline(b *rsiBot, currBid decimal.Decimal) {
	volumeToSell, funds := xrp, xbt
	price := currBid.Add(decimal.NewFromFloat64(0.00000001, 8))

	fmt.Println("SELL order placed at", price)
	b.readyToBuy = true
	b.tradesMade++

	// update funds
	xbt = xbt.Add(price.Mul(volumeToSell))
	xrp = xrp.Sub(volumeToSell)
}


// function to execute trades using historical data
func tradeOffline(b *rsiBot) {
	currRow++
	currAsk, currBid := getOfflineAsk(currRow), getOfflineBid(currRow)

	// calculating RSI using RSI algorithm
	var rsi decimal.Decimal
	rsi, b.upEma, b.downEma = getRsi(b.prevAsk, currAsk, b.upEma, b.downEma, b.tradingPeriod)
	fmt.Println("RSI", rsi, "U:", b.upEma, "D:", b.downEma)
	b.prevAsk = currAsk

	if b.readyToBuy { // check if sell order has gone trough
		fmt.Println("Current Ask", currAsk)
		if rsi.Cmp(b.overSold) == -1 && rsi.Sign() != 0 {
			buyOffline(b, currAsk)
		}
	} else {
		bound := currBid.Mul(b.stopLossMult)

		fmt.Println("Current Bid", currBid)
		fmt.Println("Stop Loss", b.stopLoss)

		if (currBid.Cmp(b.buyPrice) == 1 && currBid.Cmp(b.stopLoss) == -1) ||
			currBid.Cmp(b.buyPrice.Mul(decimal.NewFromFloat64(0.98, 8))) == -1 {
			sellOffline(b, currBid)
		} else if bound.Cmp(b.stopLoss) == 1 {
			b.stopLoss = bound
			fmt.Println("Stoploss changed to: ", b.stopLoss)
		}

	}
	b.numOfDecisions++

}
