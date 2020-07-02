package main

import (
	"github.com/luno/luno-go/decimal"
)

// function to calculate the simple moving average of a given set of data
func sma(array []decimal.Decimal) decimal.Decimal {
	sum := decimal.Zero()

	// summing elements in the given aray
	for _, val := range array {
		sum = sum.Add(val)
		if val != decimal.NewFromInt64(0) {
			listCount+=1
		}
	}
	// check to see if length of list was 0
	// in which case return 0 to avoid division by zero panic erro
	if listCount == 0 {
		return decimal.NewFromInt64(0)
	}
	return sum.Div(decimal.NewFromInt64(int64(listCount)), 8)
}


// functiont to calculate the Relative Strength Index
func rsi(array []decimal.Decimal) decimal.Decimal {
  n := len(array)

	priceUp := make([]decimal.Decimal,0)
	priceDown := make([]decimal.Decimal,0)
	//iterating through elements of array and populating priceUp/Down arrays
	for i:=0;i<n-1;i++{
		if array[i+1].Cmp(array[i]) ==1 {
			//item is over sold and therefore price is likely to rise and should buy
			diff:= array[i+1].Sub(array[i])
			frac := diff.Div(array[i],8)
			// calculating percentage rise in price
			perctangeRise := frac.Mul(decimal.NewFromInt64(100))
			priceUp = append(priceUp, perctangeRise)

		} else if array[i+1].Cmp(array[i]) == -1 {
			//item is over bought and thus price is likely to fall and should sell
			diff:= array[i].Sub(array[i+1])
			frac := diff.Div(array[i],8)
			// calculating percentage rise in price
			perctangeFall := frac.Mul(decimal.NewFromInt64(100))
			priceDown = append(priceDown, perctangeFall)
		}
	}
	// calculating average price change
	averagePriceRise := sma(priceUp)
	averagePriceFall := sma(priceDown)
	// check to see if average fall price is Zero
	// in which case return 100 to avoid non-Zero error
	comparison := decimal.NewFromInt64(1).Div(decimal.NewFromInt64(10000000),8)
	if comparison.Cmp(averagePriceFall) == 1{
		return decimal.NewFromInt64(100)

	} else {
		rs := averagePriceRise.Div(averagePriceFall,8)
		rsDen := rs.Add(decimal.NewFromInt64(1))
		// calculating rsi
		rsi1 := decimal.NewFromInt64(100).Sub(decimal.NewFromInt64(100).Div(rsDen,8))
		return rsi1
	}
}
