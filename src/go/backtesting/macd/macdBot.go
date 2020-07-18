package main

import (
  "fmt"
  "github.com/luno/luno-go/decimal"
)

type macdBot struct {
	tradingPeriodLR	int64										// How often the bot calculates a long term result
  tradingPeriodSR int64                   // How Often the bot calculates a short term result
	tradesMade	   	int64						 				// total number of trades executed
	numOfDecisions 	int64					 					// number of times the bot calculates
  data            []decimal.Decimal       // previous bids over longer trading period
	buyPrice				decimal.Decimal					// price we bought at
  currRow         int64                   // current row within excel spreadsheet
  macdValue       decimal.Decimal         // parity of short average - longer average
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


func (b* macdBot) macd() {
  LRsma := sma(b.data)
  SRsma := sma(b.data[b.tradingPeriodSR:b.tradingPeriodLR])
  macdCurr := SRsma.Sub(LRsma)
  macdPrev := b.macdValue
  if macdPrev.Cmp(decimal.Zero()) == -1 && macdCurr.Cmp(decimal.Zero()) == 1 {
    // macd line crossed over 0 line - buy
  //  fmt.Println("sma LR:", LRsma)
  //  fmt.Println("sma SR:", SRsma)
    fmt.Println("Buy")
    fmt.Println("Bought at",b.data[b.tradingPeriodLR -1] )
    b.tradesMade++
  } else if macdPrev.Cmp(decimal.Zero()) == 1 && macdCurr.Cmp(decimal.Zero()) == -1 {
    // macd line crossed under 0 line - sell
  //  fmt.Println("sma LR:", LRsma)
  //  fmt.Println("sma SR:", SRsma)
    fmt.Println("Sell")
    fmt.Println("Sold at",b.data[b.tradingPeriodLR -1] )
    b.tradesMade++
  }
  b.macdValue = macdCurr
}


func (b *macdBot) trade() {
  b.macd()
//  fmt.Println("macd:", b.macdValue)
  for i := int64(0); i < b.tradingPeriodLR - 1; i++ {
      b.data[i] = b.data[i+1]
  }
  b.currRow++
  b.data[b.tradingPeriodLR - 1] = getOfflineBid(b.currRow)

}

func (b *macdBot) initialData() {
  data := []decimal.Decimal{}
  for i:= int64(1) ; i <= b.tradingPeriodLR; i++ {
    data = append(data, getOfflineBid(i))
  }
  b.data = data
  LRsma := sma(b.data)
  SRsma := sma(b.data[b.tradingPeriodSR:b.tradingPeriodLR])

  b.macdValue = SRsma.Sub(LRsma)
}
