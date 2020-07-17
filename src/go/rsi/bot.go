package rsi

import (
	"github.com/luno/luno-go/decimal"
)

// struct for the rsiBot
type RsiBot struct {
	TradingPeriod  int64             // No of past asks used to calculate RSI
	TradesMade     int64             // total number of trades executed
	NumOfDecisions int64             // number of times the bot calculates
	StopLoss       decimal.Decimal   // variable stop loss
	StopLossMult   decimal.Decimal   // multiplier for stop loss
	OverSold       decimal.Decimal   // bound to tell the bot when to buy
	ReadyToBuy     bool              // false means ready to sell
	BuyPrice       decimal.Decimal   // stores most recent price we bought at
	UpEma					 decimal.Decimal   // exponentially smoothed Wilder's MMA for upward change
	DownEma 			 decimal.Decimal   // exponentially smoothed Wilder's MMA for downward change
	PrevAsk				 decimal.Decimal	 // the previous recorded ask price
}

// function to calculate the Relative Strength Index
func GetRsi(PrevAsk decimal.Decimal, currAsk decimal.Decimal, UpEma decimal.Decimal, DownEma decimal.Decimal, period int64) (decimal.Decimal, decimal.Decimal, decimal.Decimal) {
	//iterating through elements of array and populating priceUp/Down arrays
	var upDiff decimal.Decimal
	var downDiff decimal.Decimal

	if currAsk.Cmp(PrevAsk) == 1 {
		//item is over sold
		upDiff = currAsk.Sub(PrevAsk)
		downDiff = decimal.Zero()
	} else if currAsk.Cmp(PrevAsk) == -1 {
		//item is over bought
		upDiff = decimal.Zero()
		downDiff = PrevAsk.Sub(currAsk)
	} else {
		upDiff = decimal.Zero()
		downDiff = decimal.Zero()
	}

	priceUp := ema(UpEma, upDiff, period)
	priceDown := ema(DownEma, downDiff, period)

	// check to see if average fall price is zero to avoid div by zero error
	if priceDown.Sign() == 0 {
		if priceUp.Sign() == 0 {
			return decimal.NewFromInt64(50), priceUp, priceDown
		} else {
			return decimal.NewFromInt64(100), priceUp, priceDown
		}
	}

	rs := priceUp.Div(priceDown, 16)
	rsiDen := rs.Add(decimal.NewFromInt64(1))
	// calculating rsi
	rsi := decimal.NewFromInt64(100).Sub(decimal.NewFromInt64(100).Div(rsiDen, 16))
	return rsi, priceUp, priceDown
}
