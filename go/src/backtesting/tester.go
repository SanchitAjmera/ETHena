package main

import (
        "fmt"
        "time"
        "github.com/tealeg/xlsx"
        "github.com/luno/luno-go/decimal"
)


func getNextPrice(currRow int,fileSlice [][][]string ) decimal.Decimal {

  const priceCol = 7
  const timeCol = 0
  const sheetNum = 0

  currPrice := fileSlice[sheetNum][currRow][priceCol]
	timeStamp := fileSlice[sheetNum][currRow][timeCol]

	if currPrice == "NaN" {
		return decimal.Zero() // Zero means failed to get price
	}

	/*
  timeStampFormatted, err := time.Parse("2006-01-02 15:04:05.000", timeStamp)
  if err != nil {
    panic(err)
  }
	*/

	fmt.Println("Time: "  , timeStamp,
							"\nPrice: " , currPrice)

	currPriceDecimal, err := decimal.NewFromString(currPrice)
	if err != nil {
		panic(err)
	}
	return currPriceDecimal
}



func tester() {

  fileSlice, err := xlsx.FileToSlice("recentAPIdata.xlsx")

  if err != nil {
		panic(err)
	}
  
	bot := verySimpleBot
  sleepTime := time.Millisecond //don't need to sleep in testing
	stock := 0
	rowNum := 1
	lastPrice := getNextPrice(rowNum, fileSlice)
	startBalance := lastPrice.MulInt64(100)
	balance := startBalance
	assets:= lastPrice.MulInt64(int64(stock)).Add(balance)

	const iterations = 3 //3 days
	const minutesInDay = 1440

	for i := iterations * minutesInDay; i > 0; i--{
		rowNum++
		nextPrice := getNextPrice(rowNum, fileSlice)

		if (nextPrice.Sign() == 0) {
			fmt.Println("PRICE UNAVAILABLE")
			continue // Skip loop if price is NaN
		}

		assets = nextPrice.MulInt64(int64(stock)).Add(balance)
		fmt.Println("Balance: "  , balance,
								"\nStock: " , stock,
								"\nProfit: ", assets.Sub(startBalance),
								"\n")
	  nextTrade := bot(nextPrice, &lastPrice)
		switch {
		case nextTrade == 0:
			//do nothing
		case nextTrade > 0:
			//buy if we have enough money
			if balance.Sub(nextPrice).Sign() == 1 {
				balance = balance.Sub(nextPrice)
				stock += nextTrade
			}
		case nextTrade < 0:
			//sell if we have enough stock
			if stock + nextTrade >= 0 {
				stock += nextTrade
				balance = balance.Add(nextPrice)
			}
		default:
			panic("Unreachable")
		}
		time.Sleep(sleepTime)
	}

	//Analysis
	fmt.Println("verySimpleBot: BACKTESTING COMPLETE")
	fmt.Println("-----------------------------------")
	fmt.Println("TOTAL PROFIT: ", assets.Sub(startBalance))
	fmt.Println("PROFIT PER DAY: ", assets.Sub(startBalance).DivInt64(int64(iterations)))

}
