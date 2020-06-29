package main

import (
	// "fmt"
	"github.com/luno/luno-go/decimal"
)
/*Indicators that we can reuse. E.g. SMA*/

func sma(array []decimal.Decimal) decimal.Decimal {
	sum := decimal.Zero()
	for _, val := range array {
		sum = sum.Add(val)
	}
	return sum.DivInt64(int64(len(array)))
}
