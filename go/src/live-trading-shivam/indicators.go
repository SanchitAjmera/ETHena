package main

import (
	"github.com/luno/luno-go/decimal"
)

// function to calculate the simple moving average of a given set of data
func sma(array []decimal.Decimal) decimal.Decimal {
	sum := decimal.Zero()
	listCount := 0
	// summing elements in the given array
	for _, val := range array {
		sum = sum.Add(val)
		if val != decimal.NewFromInt64(0) {
			listCount += 1
		}
	}
	// check to see if length of list was 0
	// in which case return 0 to avoid division by zero panic error
	if listCount == 0 {
		return decimal.NewFromInt64(0)
	}
	return sum.Div(decimal.NewFromInt64(int64(listCount)), 8)
}

// function to calculate the Relative Strength Index
func rsi(array []decimal.Decimal) decimal.Decimal {
	n := len(array)

	priceUp := make([]decimal.Decimal, 0)
	priceDown := make([]decimal.Decimal, 0)
	//iterating through elements of array and populating priceUp/Down arrays
	for i := 0; i < n-1; i++ {
		if array[i+1].Cmp(array[i]) == 1 {
			//item is over sold and therefore price is likely to rise and should buy
			diff := array[i+1].Sub(array[i])
			percentageRise := diff.Div(array[i], 8)
			// calculating percentage rise in price
			priceUp = append(priceUp, percentageRise)

		} else if array[i+1].Cmp(array[i]) == -1 {
			//item is over bought and thus price is likely to fall and should sell
			diff := array[i].Sub(array[i+1])
			percentageFall := diff.Div(array[i], 8)
			// calculating percentage rise in price
			priceDown = append(priceDown, percentageFall)
		}
	}
	// calculating average price change
	averagePriceRise := sma(priceUp)
	averagePriceFall := sma(priceDown)

	// check to see if average fall price is Zero
	// in which case return 100 to avoid non-Zero error
	if averagePriceFall.Sign() == 0 {
		if averagePriceRise.Sign() == 0 {
			return decimal.NewFromInt64(50)
		} else {
			return decimal.NewFromInt64(100)
		}
	}
	rs := averagePriceRise.Div(averagePriceFall, 8)
	rsDen := rs.Add(decimal.NewFromInt64(1))
	// calculating rsi
	rsi1 := decimal.NewFromInt64(100).Sub(decimal.NewFromInt64(100).Div(rsDen, 8))
	return rsi1
}
