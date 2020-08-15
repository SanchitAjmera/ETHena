package rsi

import (
	"github.com/luno/luno-go/decimal"
)

// function to calculate the simple moving average of a given set of data
func Sma(array []decimal.Decimal, length int64) decimal.Decimal {
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
	decimalPeriod := decimal.NewFromInt64(period)
	return (oldVal.Mul(decimalPeriod.Sub(decimal.NewFromInt64(1))).Add(newData)).Div(decimalPeriod, 16)
}
