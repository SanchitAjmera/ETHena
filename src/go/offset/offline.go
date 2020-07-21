package main

import (

	"github.com/luno/luno-go/decimal"
)

var xrp decimal.Decimal
var xbt decimal.Decimal


func initialiseFunds(xbtFunds decimal.Decimal, xrpStock decimal.Decimal) {
	parseXlsx()
	xrp = xrpStock
	xbt = xbtFunds

}
/*

// function to execute buying of items
func buyOffline(b *RsiBot, currAsk decimal.Decimal) {
	currFunds := xbt
	price := currAsk.Sub(decimal.NewFromFloat64(0.00000001, 8))
	buyableStock := currFunds.Div(price, 8)
	buyableStock = buyableStock.ToScale(0)
	// checking if there are no funds available
	if currFunds.Sign() == 0 {
		fmt.Println("No funds available")
		return
	} else {
		fmt.Println("BUY order placed at", price)
		b.ReadyToBuy = false
		b.TradesMade++
		b.StopLoss = price
		b.BuyPrice = price

		// update funds
		xbt = xbt.Sub(price)
		xrp = xrp.Add(buyableStock)
	}
}

func sellOffline(b *RsiBot, currBid decimal.Decimal) {
	volumeToSell := xrp
	price := currBid.Add(decimal.NewFromFloat64(0.00000001, 8))

	fmt.Println("SELL order placed at", price)
	b.ReadyToBuy = true
	b.TradesMade++

	// update funds
	xbt = xbt.Add(price.Mul(volumeToSell))
	xrp = xrp.Sub(volumeToSell)
}


// function to execute trades using historical data
func TradeOffline(b *RsiBot) {
	currRow++
	currAsk, currBid := GetOfflineAsk(currRow), getOfflineBid(currRow)

	// calculating RSI using RSI algorithm
	var rsi decimal.Decimal
	rsi, b.UpEma, b.DownEma = GetRsi(b.PrevAsk, currAsk, b.UpEma, b.DownEma, b.TradingPeriod)
	fmt.Println("RSI", rsi, "U:", b.UpEma, "D:", b.DownEma)
	b.PrevAsk = currAsk

	printRow := currRow - 15

	f.SetCellValue("Sheet1", "B"+strconv.FormatInt(printRow, 10), rsi)
	f.SetCellValue("Sheet1", "C"+strconv.FormatInt(printRow, 10), b.ReadyToBuy)

	if b.ReadyToBuy { // check if sell order has gone trough
		f.SetCellValue("Sheet1", "A"+strconv.FormatInt(printRow, 10), currAsk)
		fmt.Println("Current Ask", currAsk)
		if rsi.Cmp(b.OverSold) == -1 && rsi.Sign() != 0 {
			buyOffline(b, currAsk)
		}
	} else {
		f.SetCellValue("Sheet1", "A"+strconv.FormatInt(printRow, 10), currBid)
		bound := currBid.Mul(b.StopLossMult)

		fmt.Println("Current Bid", currBid)
		fmt.Println("Stop Loss", b.StopLoss)

		if (currBid.Cmp(b.BuyPrice) == 1 && currBid.Cmp(b.StopLoss) == -1) ||
			currBid.Cmp(b.BuyPrice.Mul(decimal.NewFromFloat64(0.98, 8))) == -1 {
			sellOffline(b, currBid)
		} else if bound.Cmp(b.StopLoss) == 1 {
			b.StopLoss = bound
			fmt.Println("Stoploss changed to: ", b.StopLoss)
		}

	}
	b.NumOfDecisions++

}
*/
