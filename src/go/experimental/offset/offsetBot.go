package main

import (
	backtest "../../backtestingUtils"
	"fmt"
	"github.com/luno/luno-go/decimal"
	"time"
)

type offsetBot struct {
	tradingPeriod int64           // How often the bot calculates a long term result
	ema           decimal.Decimal // exponentially smoothed Wilder's MMA for upward change
	offset        decimal.Decimal // exponentially smoothed Wilder's MMA for upward changevv
	readyToBuy    bool
	StopLoss      decimal.Decimal
	StopLossMult  decimal.Decimal // multiplier for stop loss
	BuyPrice      decimal.Decimal // stores most recent price we bought at
	PrevAsk       decimal.Decimal // the previous recorded ask price
	PrevOrder     string          // stores order ID of most recent order
}

// function to calculate the simple moving average of a given set of data
func sma(array []decimal.Decimal) decimal.Decimal {
	sum := decimal.Zero()
	listCount := 0
	// summing elements in the given aray
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

func ema(oldVal decimal.Decimal, newData decimal.Decimal, period int64) decimal.Decimal {
	decimalPeriod := decimal.NewFromInt64(period)
	return (oldVal.Mul(decimalPeriod.Sub(decimal.NewFromInt64(1))).Add(newData)).Div(decimalPeriod, 16)
}

func (b *offsetBot) tradeOffline() {
	currAsk := backtest.GetOfflineAsk(currRow)
	currBid := backtest.GetOfflineBid(currRow)
	ema := ema(b.ema, currAsk, b.tradingPeriod)

	b.ema = ema
	b.PrevAsk = currAsk
	fmt.Println("EMA", b.ema, " currAsk:", currAsk, " currBid:", currBid, " CURRASK-(EMA-OFF):", currAsk.Sub(ema.Sub(b.offset)))

	if b.readyToBuy {
		if currAsk.Cmp(ema.Sub(b.offset)) == -1 {
			price := currAsk.Sub(decimal.NewFromFloat64(0.00000001, 8))
			fmt.Println("Buy     | currAsk:", price, " currRow", currRow)
			b.readyToBuy = false
			b.StopLoss = price
			b.BuyPrice = price

		}
	} else {
		if currAsk.Cmp(ema.Add(b.offset)) == 1 && currBid.Cmp(b.StopLoss.Mul(b.StopLossMult)) == 1 {
			price := currAsk.Add(decimal.NewFromFloat64(0.00000001, 8))
			b.readyToBuy = true
			fmt.Println("Sell    | currAsk:", price, "  currBid:", currBid, " currRow", currRow)
		}
		/*

			    bound := currBid.Mul(b.StopLossMult)
			    if (currBid.Cmp(b.BuyPrice) == 1 && currBid.Cmp(b.StopLoss) == -1) ||
						currBid.Cmp(b.BuyPrice.Mul(decimal.NewFromFloat64(0.98, 8))) == -1 {
			      price := currBid.Add(decimal.NewFromFloat64(0.00000001, 8))
						fmt.Println("Sell    | currBid:", price)
			      b.readyToBuy = true
			      b.funds = (b.stock.Mul(price))
			      b.stock = decimal.Zero()
					} else if bound.Cmp(b.StopLoss) == 1 {
						b.StopLoss = bound
						fmt.Println("Stoploss changed to: ", b.StopLoss)
			    }
		*/
	}

	currRow++
	//fmt.Println("curr Row: ", currRow)
}

func (b *offsetBot) tradeOnline() {

	time.Sleep(time.Minute)
	currAsk, currBid := getTicker()
	ema := ema(b.ema, currAsk, b.tradingPeriod)

	b.ema = ema
	b.PrevAsk = currAsk
	fmt.Println("EMA", b.ema, " currAsk:", currAsk, " currBid:", currBid, " CURRASK-(EMA-OFF):", currAsk.Sub(ema.Sub(b.offset)))

	if b.readyToBuy {
		if currAsk.Cmp(ema.Sub(b.offset)) == -1 {
			buy(b, currAsk)
			email("BOUGHT", currAsk, currAsk.Sub(ema.Sub(b.offset)))
		}
	} else {
		if currAsk.Cmp(ema.Add(b.offset)) == 1 && currBid.Cmp(b.StopLoss.Mul(b.StopLossMult)) == 1 {
			sell(b, currBid)
			email("SOLD", currBid, (ema.Add(b.offset)).Sub(currAsk))
		}
	}
}
