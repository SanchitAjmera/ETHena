package main

import (
        "fmt"
        "github.com/tealeg/xlsx"
        "strconv"
)


func getNextPrice(currRow int,fileSlice [][][]string ) (float64, float64) {

  currPrice := fileSlice[sheetNum][currRow][priceCol]
	timeStamp := fileSlice[sheetNum][currRow][timeCol]

	if currPrice == "NaN" {
		return float64(0), float64(0) // Zero means failed to get price
	}

	currPriceDecimal, err := strconv.ParseFloat(currPrice, 8)
	if err != nil {
		panic(err)
	}

  timeStampFloat, err := strconv.ParseFloat(timeStamp, 8)
  if err != nil {
		panic(err)
	}

	return currPriceDecimal, timeStampFloat
}



func tester() {
  fileSlice, err := xlsx.FileToSlice("recentAPIdata.xlsx")

  if err != nil {
		panic(err)
	}

  inventory := [][]float64{}

  for i := 1; i<720; i++ {

    trader1 := &state_t{funds : float64(100000), assets: float64(0),
      inventory : inventory, historicalData : fileSlice, currentDay : i}

    trader1.metrics.tickerTime = 0
    trader1.metrics.dataCacheLength = trader1.currentDay
    trader1.metrics.offset = float64(100)


    day := 0

    for j := 0 ; j < 1440; j++ {
      day +=trader1.metrics.dataCacheLength
  //    fmt.Println(trader1.currentDay)
      if day > 1440{
        break
      }
      SMEBot(trader1, false)

    }

/*    fmt.Println(".")
    fmt.Println(".")
    fmt.Println(".")
    fmt.Println(".")
*/
    lastPrice, _ := getNextPrice(trader1.currentDay, fileSlice)
    sell(trader1, lastPrice, false)
/*
    fmt.Println(".")
    fmt.Println(".")
    fmt.Println(".")
    fmt.Println(".")
*/  fmt.Println("----------------------------------------------------------------")
    fmt.Println(".")
    fmt.Println("                   currentDay:  ", trader1.currentDay)
    fmt.Println("                   timePast:    ", trader1.metrics.dataCacheLength)
    fmt.Println("                   Initial Funds: £ 100000")
    fmt.Println("                   Final Funds:   £", trader1.funds)
    fmt.Println(".")
    fmt.Println(".")

    }

  }


/*
func tester2() {

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


*/
