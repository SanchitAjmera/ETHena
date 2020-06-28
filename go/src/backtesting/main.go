package main


import (
        "fmt"
        "github.com/luno/luno-go/decimal"
				"github.com/tealeg/xlsx"
)

func main () {
	fileSlice, err := xlsx.FileToSlice("recentAPIdata.xlsx")
	if err != nil {
		panic(err)
	}

	const priceCol = 7
	const timeCol = 0
	const sheetNum = 0

	var getNextPrice ticker = func (currRow int) decimal.Decimal {
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

	tester(getNextPrice)
}
