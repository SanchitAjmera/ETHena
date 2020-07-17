package main

import (
	"fmt"
	"github.com/luno/luno-go/decimal"
)

// struct for portfolio which includes features common to every bot
type portfolio struct {
	funds         decimal.Decimal  // money bot has to trade with
	stock         decimal.Decimal  // number of bought and not yet sold items
	TradingPeriod int64 					 // How often the bot calculates a result
	currRow       int64						 // current row within excel spreadsheet
	TradesMade		int							 // total number of trades executed
	StopLoss		 	 decimal.Decimal  // variable stop loss
	StopLossMult   decimal.Decimal  // multiplier for stop loss
	ReadyToBuy		 bool
}

// struct for smaBot
type smaBot struct {
	pf             *portfolio			 // portfolio of bot
	offset         decimal.Decimal // Offset size
	NumOfDecisions int64           // number of times the bot calculates
}

// struct for the rsiBot
type rsiBot struct {
	pf 						 *portfolio			 // portfolio of bot
	NumOfDecisions int64					 // number of times the bot calculates
	OverSold			 int64					 // bound to tell if the item is over sold
	overBought		 int64					 // bound to tell if the item is over bought
	prevBid				 decimal.Decimal
}

type macdBot struct {
	pf 						 *portfolio
	NumOfDecisions int64
	SRtradePeriod  int64
	LRtradePeriod  int64
	diffs					 []decimal.Decimal
	prevDiff			 decimal.Decimal
}

// function to execute buying of items
func buy(pf *portfolio, stock decimal.Decimal, price decimal.Decimal) {
	currFunds := pf.funds
	// checking if there are enough funds to buy the given amount of stock
	if currFunds.Cmp(stock.Mul(price)) == -1 {
		fmt.Println("Cannot afford to buy ", stock, " stock at ", price, " price")
	} else {
		// updating portfolio of trader bot
		// fmt.Println("Bought ", stock, " stock at ", price, " price")
		pf.funds = pf.funds.Sub(stock.Mul(price))
		pf.stock = pf.stock.Add(stock)
		pf.TradesMade++
		pf.ReadyToBuy = false
		pf.StopLoss = price
		//fmt.Println("buy")
	  //printPortFolio(pf)
		// sets new stop loss to new price if price > current stop loss
		//if pf.StopLoss.Cmp(price) == -1 {
		//}
		// fmt.Println("Current funds: ",pf.funds,"\n")
	}
}

// function to execute selling of items
func sell(pf *portfolio, stock decimal.Decimal, price decimal.Decimal) {
	currStock := pf.stock
	// stops sell if price < stoploss * multiplier
	//checking if we have enough stock to sell
	if currStock.Cmp(stock) == -1{
		fmt.Println("Not enough stock to sell ", stock, " stock at ", price, " price")
	} else {
		// updating portfolio of trader bot
		// fmt.Println("Sold ", stock, " stock at ", price, " price")
		pf.funds = pf.funds.Add(stock.Mul(price))
		pf.stock = pf.stock.Sub(stock)
		pf.TradesMade++
		pf.ReadyToBuy = true
		//fmt.Println("sell")
	  //printPortFolio(pf)
		// fmt.Println("Current funds: ",pf.funds,"\n")
	}
}

// function to execute trades of the SMA bot
func (b *smaBot) tradeSMA() {
	// initialising variables
	pastBids := make([]decimal.Decimal, b.pf.TradingPeriod)
	var currBid decimal.Decimal
	var currAsk decimal.Decimal
	var i int64 = 0
	// getting live bid price
	currBid = getBid(b.pf.currRow)
	// getting live ask price
	currAsk = getAsk(b.pf.currRow)
	// populating past bids with bids from the last trading period
	for i < b.pf.TradingPeriod {
		pastBids[i] = getBid(b.pf.currRow - i)
		i++
	}
	// calculating mean using the SMA algorithm
	buyableStock := b.pf.funds.Div(currAsk, 8)
	mean := sma(pastBids)

	if currBid.Cmp(mean.Add(b.offset)) == 1 && b.pf.stock.Sign() != 0{
		// selling if bid is greater than mean plus the offset
		sell(b.pf, b.pf.stock, currBid)
	} else if currBid.Cmp(mean.Sub(b.offset)) == -1 {
		// buying if bid is less than mean minus offset
		buy(b.pf, buyableStock, currAsk)
	}
	// incrementing current row by the trade period
	b.pf.currRow += b.pf.TradingPeriod
}

// function to execute trades using the RSI bot
func (b *rsiBot) tradeRSI() {

	pastAsks:= make([]decimal.Decimal, b.pf.TradingPeriod)

	for i := 0;  i < int(b.pf.TradingPeriod); i++ {
		pastAsks[i] = getAsk(b.pf.currRow - int64(i))
	}

	rsi := rsi(pastAsks)

	if b.pf.ReadyToBuy {
		if rsi.Cmp(decimal.NewFromInt64(b.OverSold)) == -1 {

			currAsk := getAsk(b.pf.currRow)
			buyableStock := b.pf.funds.Div(currAsk, 8)
			buy(b.pf, buyableStock, currAsk)
			price := currAsk.Mul(decimal.NewFromFloat64(0.99999,8))
			b.pf.StopLoss = price
			fmt.Println("BUYING AT £", price)
			printPortFolio(b.pf)
		}
	} else {
		currBid := getCurrBid()
		bound := currBid.Mul(b.pf.StopLossMult)

		if b.prevBid.Cmp(b.pf.StopLoss) == 1 && currBid.Cmp(b.pf.StopLoss) == -1 {
			price := currBid.Mul(decimal.NewFromFloat64(1.00001, 8))
			sell(b.pf, b.pf.stock, currBid)
			fmt.Println("SELLING at ", price)
			b.pf.ReadyToBuy = true
			printPortFolio(b.pf)

		} else if bound.Cmp(b.pf.StopLoss) == 1 {
			b.pf.StopLoss = bound
			fmt.Println("Stoploss: ",b.pf.StopLoss)
		}
	}
	b.NumOfDecisions+= b.pf.TradingPeriod
}



func (b *macdBot) tradeMACD(){
	currAsk := getAsk(b.pf.currRow)
	currBid := getBid(b.pf.currRow)
	pastBidsSR := make([]decimal.Decimal, b.SRtradePeriod)
	pastBidsLR := make([]decimal.Decimal, b.LRtradePeriod)

	for i := 0;  i < int(b.SRtradePeriod); i++ {
		pastBidsSR[i] = getBid(b.pf.currRow - int64(i))
	}

	for i := 0;  i < int(b.LRtradePeriod); i++ {
		pastBidsLR[i] = getBid(b.pf.currRow - int64(i))
	}

	smaLR := sma(pastBidsLR)
	smaSR := sma(pastBidsSR)

	diff := smaLR.Sub(smaSR)
	b.diffs = append(b.diffs, diff)

	buyableStock := b.pf.funds.Div(currAsk, 8)

	upperBound, err := decimal.NewFromString("0.01")
	if err != nil { panic(err)}
	lowerBound, err := decimal.NewFromString("-0.01")
	if err != nil { panic(err)}

	if diff.Cmp(upperBound) == -1 && diff.Cmp(lowerBound) == 1 {

		if b.prevDiff.Cmp(decimal.Zero()) == 1 {
			// coming from smaLR > smaSR
			buy(b.pf, buyableStock, currAsk)

		} else if b.prevDiff.Cmp(decimal.Zero()) == -1 {
			// coming from smaLR > smaSR
			sell(b.pf, b.pf.stock, currBid)
		}
		//fmt.Println("SR :", smaSR)
		//fmt.Println("LR :", smaLR)
		//fmt.Println("\n timestamp: ", b.pf.currRow, "\n")
		//fmt.Println("crossOver")
	}

	b.prevDiff = diff
	b.pf.currRow += b.pf.TradingPeriod

}
