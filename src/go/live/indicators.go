package main

import (
	"github.com/luno/luno-go/decimal"
)

// function to calculate the simple moving average of a given set of data
func sma(array []decimal.Decimal, length int64) decimal.Decimal {
	// check to see if length of list was 0
	// in which case return 0 to avoid division by zero panic error
	if length == 0 {
		return decimal.Zero()
	}

	sum := decimal.Zero()
	// summing elements in the given array
	for _, val := range array {
		sum = sum.Add(val)
	}
	return sum.Div(decimal.NewFromInt64(length), 16)
}

//function to calculate exponentially smoothed moving average
func ema(oldVal decimal.Decimal, newData decimal.Decimal, period int64) decimal.Decimal {
	return ((oldVal.Mul(period.Sub(decimal.NewFromInt64(1))).Add(newData)).Div(period, 16)
}

// function to calculate the Relative Strength Index
func getRsi(prevAsk decimal.Decimal, currAsk decimal.Decimal, upEma decimal.Decimal, downEma decimal.Decimal, period int64) (decimal.Decimal, decimal.Decimal, decimal.Decimal) {
	//iterating through elements of array and populating priceUp/Down arrays
	var upDiff decimal.Decimal
	var downDiff decimal.Decimal

	if currAsk.Cmp(prevAsk) == 1 {
		//item is over sold
		upDiff = currAsk.Sub(prevAsk)
		downDiff = decimal.Zero()
	} else if currAsk.Cmp(prevAsk) == -1 {
		//item is over bought
		upDiff = decimal.Zero()
		downDiff = prevAsk.Sub(currAsk)
	} else {
		upDiff = decimal.Zero()
		downDiff = decimal.Zero()
	}

	priceUp := ema(upEma, upDiff, period)
	priceDown := ema(downEma, downDiff, period)

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
