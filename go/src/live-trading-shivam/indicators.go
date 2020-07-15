package main

import (
	"github.com/luno/luno-go/decimal"
)

// function to calculate the simple moving average of a given set of data
func sma(array []decimal.Decimal, length int64) decimal.Decimal {
	// check to see if length of list was 0
	// in which case return 0 to avoid division by zero panic error
	if length == 0 {
		return decimal.NewFromInt64(0)
	}

	sum := decimal.Zero()
	// summing elements in the given array
	for _, val := range array {
		sum = sum.Add(val)
	}
	return sum.Div(decimal.NewFromInt64(length), 16)
}

//function to calculate exponentially smoothed moving average
func ema(oldVal decimal.Decimal, newVal decimal.Decimal, period int64) decimal.Decimal {
	alpha := decimal.NewFromFloat64(1.0/float64(period), 16)
	return alpha.Mul(newVal).Add(oldVal.Mul(decimal.NewFromInt64(1).Sub(alpha)))
}

// function to calculate the Relative Strength Index
func getRsi(prevAsk decimal.Decimal, currAsk decimal.Decimal, upEma decimal.Decimal, downEma decimal.Decimal, period int64) (decimal.Decimal, decimal.Decimal, decimal.Decimal) {
	priceUp := decimal.Zero()
	priceDown := decimal.Zero()
	//iterating through elements of array and populating priceUp/Down arrays
	if currAsk.Cmp(prevAsk) == 1 {
		//item is over sold
		diff := currAsk.Sub(prevAsk)
		priceUp = ema(upEma, diff, period)
		priceDown = ema(downEma, priceDown, period)
	} else if currAsk.Cmp(prevAsk) == -1 {
		//item is over bought
		diff := prevAsk.Sub(currAsk)
		priceUp = ema(downEma, priceUp, period)
		priceDown = ema(downEma, diff, period)
	}
	// check to see if average fall price is zero to avoid div by zero error
	if priceDown.Sign() == 0 {
		if priceUp.Sign() == 0 {
			return decimal.NewFromInt64(50), priceUp, priceDown
		} else {
			return decimal.NewFromInt64(100), priceUp, priceDown
		}
	}
	rs := priceUp.Div(priceDown, 16)
	rsDen := rs.Add(decimal.NewFromInt64(1))
	// calculating rsi
	rsi := decimal.NewFromInt64(100).Sub(decimal.NewFromInt64(100).Div(rsDen, 16))
	return rsi, priceUp, priceDown
}
