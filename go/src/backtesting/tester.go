package main

import (
        "fmt"
        "time"
        "github.com/tealeg/xlsx"
        "strconv"
)


func getNextPrice(currRow int,fileSlice [][][]string ) float64 {

  currPrice := fileSlice[sheetNum][currRow][priceCol]
	timeStamp := fileSlice[sheetNum][currRow][timeCol]

	if currPrice == "NaN" {
		return float64(0) // Zero means failed to get price
	}

	/*
  timeStampFormatted, err := time.Parse("2006-01-02 15:04:05.000", timeStamp)
  if err != nil {
    panic(err)
  }
	*/

	fmt.Println("Time: "  , timeStamp,
							"\nPrice: " , currPrice)

	currPriceDecimal, err := strconv.ParseFloat(currPrice, 8)
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
	startBalance := float64(1000000)
	balance := startBalance
	assets:= lastPrice * float64(stock) + balance

	const iterations = 3 //3 days
	const minutesInDay = 1440

	for i := iterations * minutesInDay; i > 0; i--{
		rowNum++
		nextPrice := getNextPrice(rowNum, fileSlice)

		if (nextPrice == float64(0)) {
			fmt.Println("PRICE UNAVAILABLE")
			continue // Skip loop if price is NaN
		}

		assets = nextPrice * float64(stock) + balance
		fmt.Println("Balance: "  , balance,
								"\nStock: " , stock,
								"\nProfit: ", assets - startBalance,
								"\n")
	  nextTrade := bot(nextPrice, &lastPrice)
		switch {
		case nextTrade == 0:
			//do nothing
		case nextTrade > 0:
			//buy if we have enough money
			if balance - nextPrice > 0 {
				balance = balance - nextPrice
				stock += int(nextTrade)
			}
		case nextTrade < 0:
			//sell if we have enough stock
			if float64(stock) + nextTrade >= 0 {
				stock += int(nextTrade)
				balance = balance - nextPrice
			}
		default:
			panic("Unreachable")
		}
		time.Sleep(sleepTime)
	}

	//Analysis
	fmt.Println("verySimpleBot: BACKTESTING COMPLETE")
	fmt.Println("-----------------------------------")
	fmt.Println("TOTAL PROFIT: ", assets - startBalance)
	fmt.Println("PROFIT PER DAY: ", assets - startBalance / float64(iterations))

}
