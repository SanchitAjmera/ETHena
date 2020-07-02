package main

import (
	"github.com/luno/luno-go/decimal"
)
/*Indicators that we can reuse. E.g. SMA*/

func sma(array []decimal.Decimal) decimal.Decimal {
	sum := decimal.Zero()
	listCount := 0

	for _, val := range array {
		sum = sum.Add(val)
		if val != decimal.NewFromInt64(0) {
			listCount+=1
		}
	}

	if listCount == 0 {
		return decimal.NewFromInt64(0)
	}

	return sum.Div(decimal.NewFromInt64(int64(listCount)), 8)
}


func rsi(array []decimal.Decimal) decimal.Decimal {
  n := len(array)

	priceUp := make([]decimal.Decimal,0)
	priceDown := make([]decimal.Decimal,0)

	for i:=0;i<n-1;i++{

		if array[i+1].Cmp(array[i]) ==1 {
			//price goes up
			diff:= array[i+1].Sub(array[i])
			frac := diff.Div(array[i],8)
			perctangeRise := frac.Mul(decimal.NewFromInt64(100))
			priceUp = append(priceUp, perctangeRise)

		} else if array[i+1].Cmp(array[i]) == -1 {
			//price goes down
			diff:= array[i].Sub(array[i+1])
			frac := diff.Div(array[i],8)
			perctangeFall := frac.Mul(decimal.NewFromInt64(100))
			priceDown = append(priceDown, perctangeFall)
		}
	}

	averagePriceRise := sma(priceUp)
	averagePriceFall := sma(priceDown)

	comparison := decimal.NewFromInt64(1).Div(decimal.NewFromInt64(10000000),8)
	if comparison.Cmp(averagePriceFall) == 1{
		return decimal.NewFromInt64(100)

	} else{
		rs := averagePriceRise.Div(averagePriceFall,8)
		rsDen := rs.Add(decimal.NewFromInt64(1))

		rsi1 := decimal.NewFromInt64(100).Sub(decimal.NewFromInt64(100).Div(rsDen,8))


		return rsi1
	}
}
