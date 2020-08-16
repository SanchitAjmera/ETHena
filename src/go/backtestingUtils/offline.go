package backtestingUtils

import (
	live "../Utils"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/luno/luno-go/decimal"
	"strconv"
)

var xrp decimal.Decimal
var xbt decimal.Decimal
var currRow int64
var f *excelize.File

func InitialiseFunds(xbtFunds decimal.Decimal, xrpStock decimal.Decimal) {
	parseXlsx()
	xrp = xrpStock
	xbt = xbtFunds
	currRow = 0
	f = excelize.NewFile()
	f.SetCellValue("Sheet1", "A1", "Curr Price")
	f.SetCellValue("Sheet1", "B1", "RSI")
	f.SetCellValue("Sheet1", "C1", "Mode (ReadyToBuy?)")
}

// function to execute buying of items
func buyOffline(b *live.RsiBot, currAsk decimal.Decimal) {
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

func sellOffline(b *live.RsiBot, currBid decimal.Decimal) {
	volumeToSell := xrp
	price := currBid.Add(decimal.NewFromFloat64(0.00000001, 8))

	fmt.Println("SELL order placed at", price)
	b.ReadyToBuy = true
	b.TradesMade++

	// update funds
	xbt = xbt.Add(price.Mul(volumeToSell))
	xrp = xrp.Sub(volumeToSell)
}

// TradeOffline function to execute trades using historical data
func TradeOffline(b *live.RsiBot) {

	currAsk, currBid := GetOfflineAsk(currRow+b.LongestTradingPeriod), GetOfflineBid(currRow+b.LongestTradingPeriod)

	rsiweighting := int(b.BotString[0])
	MACDweighting := int(b.BotString[1])
	Candlestickweighting := int(b.BotString[2])
	Offsetweighting := int(b.BotString[3])

	b.PastAsks = b.PastAsks[1:]
	b.PastAsks = append(b.PastAsks, currAsk)
	// calculating RSI using RSI algorithm
	var rsi decimal.Decimal
	scores := []decimal.Decimal{}
	prevema := live.Sma(b.PastAsks[b.LongestTradingPeriod-b.OffsetTraingPeriod : b.LongestTradingPeriod-1])
	rsi, b.UpEma, b.DownEma = live.GetRsi(b.PrevAsk, currAsk, b.UpEma, b.DownEma, b.RSITradingPeriod)
	if rsiweighting != '0' {
		rsiScore := decimal.NewFromInt64(100).Sub(rsi)
		for i := 0; i < rsiweighting; i++ {
			scores = append(scores, rsiScore)
		}
	}

	b.PrevAsk = currAsk

	if MACDweighting != '0' {
		b.MACDlongperiodavg = live.Sma(b.PastAsks[b.LongestTradingPeriod-b.MACDTradingPeriodLR:])
		b.MACDshortperiodavg = live.Sma(b.PastAsks[b.LongestTradingPeriod-b.MACDTradingPeriodSR:])
		currdifference := b.MACDshortperiodavg.Sub(b.MACDlongperiodavg)
		macdScore := decimal.NewFromInt64(100).Sub(currdifference.Div(decimal.NewFromFloat64(0.000001, 16), 16))
		for i := 0; i < MACDweighting; i++ {
			scores = append(scores, macdScore)
		}
	}

	if Candlestickweighting != '0' {
		if live.Rev123(b.Stack[b.LongestTradingPeriod-3], b.Stack[b.LongestTradingPeriod-2], b.Stack[b.LongestTradingPeriod-1]) || live.Hammer(b.Stack[b.LongestTradingPeriod-1]) || live.InverseHammer(b.Stack[b.LongestTradingPeriod-1]) || live.WhiteSlaves(b.Stack[b.LongestTradingPeriod-3], b.Stack[b.LongestTradingPeriod-2], b.Stack[b.LongestTradingPeriod-1]) || live.MorningStar(b.Stack[b.LongestTradingPeriod-3], b.Stack[b.LongestTradingPeriod-2], b.Stack[b.LongestTradingPeriod-1]) {
			candlestickscore := decimal.NewFromInt64(100)
			for i := 0; i < Candlestickweighting; i++ {
				scores = append(scores, candlestickscore)
			}
		}
	}
	if Offsetweighting != '0' {
		ema := live.Ema(prevema, currAsk, b.OffsetTraingPeriod)
		if currAsk.Cmp(ema.Sub(b.Offset)) == -1 {
			offsetscore := decimal.NewFromInt64(100)
			for i := 0; i < Offsetweighting; i++ {
				scores = append(scores, offsetscore)
			}
		}
	}

	currRow++

	averageScore := live.Sma(scores)

	printRow := currRow - 15

	f.SetCellValue("Sheet1", "B"+strconv.FormatInt(printRow, 10), rsi)
	f.SetCellValue("Sheet1", "C"+strconv.FormatInt(printRow, 10), b.ReadyToBuy)

	if b.ReadyToBuy { // check if sell order has gone trough
		f.SetCellValue("Sheet1", "A"+strconv.FormatInt(printRow, 10), currAsk)
		fmt.Println("Current Ask", currAsk)
		if averageScore.Cmp(decimal.NewFromInt64(80)) == 1 {
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
