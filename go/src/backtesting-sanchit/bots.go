package main

import (
	"fmt"
	"github.com/luno/luno-go/decimal"
)

// struct for portfolio which includes features common to every bot
type portfolio struct {
	funds         decimal.Decimal  // money bot has to trade with
	stock         decimal.Decimal  // number of bought and not yet sold items
	tradingPeriod int64 					 // How often the bot calculates a result
	currRow       int64						 // current row within excel spreadsheet
	tradesMade		int							 // total number of trades executed
}

// struct for smaBot
type smaBot struct {
	pf             *portfolio			 // portfolio of bot
	offset         decimal.Decimal // Offset size
	numOfDecisions int64           // number of times the bot calculates
}

// struct for the rsiBot
type rsiBot struct {
	pf 						 *portfolio			 // portfolio of bot
	numOfDecisions int64					 // number of times the bot calculates
	overSold			 int64					 // bound to tell if the item is over sold
	overBought		 int64					 // bound to tell if the item is over bought
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
		pf.tradesMade++
		// fmt.Println("Current funds: ",pf.funds,"\n")
	}
}

// function to execute selling of items
func sell(pf *portfolio, stock decimal.Decimal, price decimal.Decimal) {
	currStock := pf.stock
	//checking if we have enough stock to sell
	if currStock.Cmp(stock) == -1{
		fmt.Println("Not enough stock to sell ", stock, " stock at ", price, " price")
	} else {
		// updating portfolio of trader bot
		// fmt.Println("Sold ", stock, " stock at ", price, " price")
		pf.funds = pf.funds.Add(stock.Mul(price))
		pf.stock = pf.stock.Sub(stock)
		pf.tradesMade++
		// fmt.Println("Current funds: ",pf.funds,"\n")
	}
}

// function to execute trades of the SMA bot
func (b *smaBot) tradeSMA() {
	// initialising variables
	pastBids := make([]decimal.Decimal, b.pf.tradingPeriod)
	var currBid decimal.Decimal
	var currAsk decimal.Decimal
	var i int64 = 0
	// getting live bid price
	currBid = getBid(b.pf.currRow)
	// getting live ask price
	currAsk = getAsk(b.pf.currRow)
	// populating past bids with bids from the last trading period
	for i < b.pf.tradingPeriod {
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
	b.pf.currRow += b.pf.tradingPeriod
}

// function to execute trades using the RSI bot
func (b *rsiBot) tradeRSI() {
	// initiailising metrics and empty array
	var currBid decimal.Decimal
	var currAsk decimal.Decimal
	var i int64 = 0
	pastBids := make([]decimal.Decimal, b.pf.tradingPeriod)
	overSold := b.overSold
	overBought := b.overBought
	// getting live bid price
	currBid = getBid(b.pf.currRow)
	// getting live ask price
	currAsk = getAsk(b.pf.currRow)
	// populating past bids array with bids from the previos trading period
	for i < b.pf.tradingPeriod {
		pastBids[i] = getBid(b.pf.currRow - i)
		i++
	}
	// calculating RSI usig RSI algorithm
	buyableStock := b.pf.funds.Div(currAsk, 8)
	rsi := rsi(pastBids)

	if rsi.Cmp(decimal.NewFromInt64(overSold)) == -1 {
		// buying stock if rsi is less than overSold bound
		buy(b.pf, buyableStock, currAsk)
	} else if rsi.Cmp(decimal.NewFromInt64(overBought)) == 1 {
		// selling stokc if rsi is greater than overBought bound
		sell(b.pf, b.pf.stock, currBid)
	}
	// incrementing current row by the trading period
	b.pf.currRow += b.pf.tradingPeriod
}
