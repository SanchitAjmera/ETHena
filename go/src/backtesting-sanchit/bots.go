package main

import (
	"fmt"

	"github.com/luno/luno-go/decimal"
)

type portfolio struct { //Features common to every bot
	funds         decimal.Decimal
	stock         decimal.Decimal
	tradingPeriod int64 //Trading period in minutes
	currRow       int64
	tradesMade		int
}

type smaBot struct { //Wagwan this is that
	pf             *portfolio
	offset         decimal.Decimal //Offset size
	numOfDecisions int64           //Length of short moving average as multiple of period
}


type rsiBot struct {
	pf 						 *portfolio
	numOfDecisions int64
	overSold			 int64
	overBought		 int64
}


func buy(pf *portfolio, stock decimal.Decimal, price decimal.Decimal) {
	currFunds := pf.funds
	if currFunds.Cmp(stock.Mul(price)) == -1 {
		fmt.Println("Cannot afford to buy ", stock, " stock at ", price, " price")
	} else {
		// fmt.Println("Bought ", stock, " stock at ", price, " price")
		pf.funds = pf.funds.Sub(stock.Mul(price))
		pf.stock = pf.stock.Add(stock)
		pf.tradesMade++
		// fmt.Println("Current funds: ",pf.funds,"\n")
	}
}

func sell(pf *portfolio, stock decimal.Decimal, price decimal.Decimal) {
	currStock := pf.stock
	if currStock.Cmp(stock) == -1{
		fmt.Println("Not enough stock to sell ", stock, " stock at ", price, " price")
	} else {
		// fmt.Println("Sold ", stock, " stock at ", price, " price")
		pf.funds = pf.funds.Add(stock.Mul(price))
		pf.stock = pf.stock.Sub(stock)
		pf.tradesMade++
		// fmt.Println("Current funds: ",pf.funds,"\n")
	}
}


func (b *smaBot) tradeSMA() {
	pastBids := make([]decimal.Decimal, b.pf.tradingPeriod)
	var currBid decimal.Decimal
	var currAsk decimal.Decimal

	var i int64 = 0

	currBid = getBid(b.pf.currRow)

	for i < b.pf.tradingPeriod {
		pastBids[i] = getBid(b.pf.currRow - i)
		i++
	}

	currAsk = getAsk(b.pf.currRow)


	buyableStock := b.pf.funds.Div(currAsk, 8)
	mean := sma(pastBids)

	if currBid.Cmp(mean.Add(b.offset)) == 1 && b.pf.stock.Sign() != 0{
		sell(b.pf, b.pf.stock, currBid)
	} else if currBid.Cmp(mean.Sub(b.offset)) == -1 {
		buy(b.pf, buyableStock, currAsk)
	}

	b.pf.currRow += b.pf.tradingPeriod
}


func (b *rsiBot) tradeRSI() {
	pastBids := make([]decimal.Decimal, b.pf.tradingPeriod)
	var currBid decimal.Decimal
	var currAsk decimal.Decimal
	overSold := b.overSold
	overBought := b.overBought

	var i int64 = 0

	currBid = getBid(b.pf.currRow)

	for i < b.pf.tradingPeriod {
		pastBids[i] = getBid(b.pf.currRow - i)
		i++
	}

	currAsk = getAsk(b.pf.currRow)


	buyableStock := b.pf.funds.Div(currAsk, 8)
	rsi := rsi(pastBids)

	if rsi.Cmp(decimal.NewFromInt64(overSold)) == -1 {
		buy(b.pf, buyableStock, currAsk)
	} else if rsi.Cmp(decimal.NewFromInt64(overBought)) == 1 {
		sell(b.pf, b.pf.stock, currBid)
	}

	b.pf.currRow += b.pf.tradingPeriod
}
