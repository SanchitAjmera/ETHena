package main

import (
  "fmt"
  "github.com/luno/luno-go/decimal"
)

type offsetBot struct {
	tradingPeriod	int64										// How often the bot calculates a long term result
  ema				  	  decimal.Decimal         // exponentially smoothed Wilder's MMA for upward change
  currRow         int64                   // current row within excel spreadsheet
  offset				  decimal.Decimal         // exponentially smoothed Wilder's MMA for upward changevv
  readyToBuy      bool
}

// function to calculate the simple moving average of a given set of data
func sma(array []decimal.Decimal) decimal.Decimal {
	sum := decimal.Zero()
	listCount := 0
	// summing elements in the given aray
	for _, val := range array {
		sum = sum.Add(val)
		if val != decimal.NewFromInt64(0) {
			listCount+=1
		}
	}
	// check to see if length of list was 0
	// in which case return 0 to avoid division by zero panic error
	if listCount == 0 {
		return decimal.NewFromInt64(0)
	}
	return sum.Div(decimal.NewFromInt64(int64(listCount)), 8)
}

func ema(oldVal decimal.Decimal, newData decimal.Decimal, period int64) decimal.Decimal {
	decimalPeriod := decimal.NewFromInt64(period)
	return (oldVal.Mul(decimalPeriod.Sub(decimal.NewFromInt64(1))).Add(newData)).Div(decimalPeriod, 16)
}

func (b *offsetBot) trade(){
  currAsk := getOfflineAsk(b.currRow)
  ema := ema(b.ema, currAsk, b.tradingPeriod)

  b.ema = ema

  if currAsk.Cmp(ema.Add(b.offset)) == 1  && !b.readyToBuy {
    fmt.Println("sell   | currAsk:", currAsk)
    b.readyToBuy = true

  } else if currAsk.Cmp(ema.Sub(b.offset)) == -1 && b.readyToBuy {
    fmt.Println("Buy    | currAsk:", currAsk)
    b.readyToBuy = false
  }
  b.currRow++
  //fmt.Println("curr Row: ", b.currRow)
}
