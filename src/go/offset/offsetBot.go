package main

import (
  "fmt"
  "github.com/luno/luno-go/decimal"
)

type offsetBot struct {
	tradingPeriod	int64										// How often the bot calculates a long term result
  emaAsk				  	  decimal.Decimal         // exponentially smoothed Wilder's MMA for upward change
  emaBid				  	  decimal.Decimal         // exponentially smoothed Wilder's MMA for upward change
  currRow         int64                   // current row within excel spreadsheet
  offset				  decimal.Decimal         // exponentially smoothed Wilder's MMA for upward changevv
  readyToBuy      bool
  StopLoss        decimal.Decimal
	StopLossMult    decimal.Decimal   // multiplier for stop loss
  BuyPrice        decimal.Decimal   // stores most recent price we bought at
  PrevAsk				 decimal.Decimal	 // the previous recorded ask price
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
  currBid := getOfflineBid(b.currRow)
  emaAsk := ema(b.emaAsk, currAsk, b.tradingPeriod)
  emaBid := ema(b.emaBid, currBid, b.tradingPeriod)

  b.emaAsk = emaAsk
  b.emaBid = emaBid
  b.PrevAsk = currAsk

  if b.readyToBuy {
    if currAsk.Cmp(emaAsk.Sub(b.offset)) == -1 {
      price := currAsk.Sub(decimal.NewFromFloat64(0.00000001, 8))
      fmt.Println("Buy     | currAsk:", price)
      b.readyToBuy = false
      b.StopLoss = price
      b.BuyPrice = price

    }
  } else {
    if currAsk.Cmp(emaAsk.Add(b.offset)) == 1  && currBid.Cmp(b.StopLoss.Mul(b.StopLossMult)) == 1 {
      price := currAsk.Add(decimal.NewFromFloat64(0.00000001, 8))
      fmt.Println("Sell    | currAsk:", price, "  currBid:", currBid)
      b.readyToBuy = true
    }

/*
    bound := currBid.Mul(b.StopLossMult)
    if (currBid.Cmp(b.BuyPrice) == 1 && currBid.Cmp(b.StopLoss) == -1) ||
			currBid.Cmp(b.BuyPrice.Mul(decimal.NewFromFloat64(0.98, 8))) == -1 {
      price := currBid.Add(decimal.NewFromFloat64(0.00000001, 8))
			fmt.Println("Sell    | currBid:", price)
      b.readyToBuy = true
		} else if bound.Cmp(b.StopLoss) == 1 {
			b.StopLoss = bound
			fmt.Println("Stoploss changed to: ", b.StopLoss)
    }
  */
  }

  b.currRow++
  //fmt.Println("curr Row: ", b.currRow)
}
