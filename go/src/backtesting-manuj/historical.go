package main

import (
  "github.com/tealeg/xlsx"
  "strconv"
)

/*
TODO:
- change colNum far ask and bid
*/

var historicalData [][][]String

// function to process the csv file and return a 3d array of strings
// historicalData is of the form: [sheetNum][rowNum][colNum]
func parse_xlsx (filename String) {
  fileSlice, err := xlsx.FileToSlice("recentAPIdata.xlsx")

  if err != nil {
    panic(err)
  }

  historicalData := fileSlice
}

func getBid (currRow int) decimal.Decimal {

  currPrice := fileSlice[0][currRow][7] //Change this
	currPriceDecimal, err := decimal.NewFromString(currPrice)

	if err != nil {
		panic(err)
	}

	return currPriceDecimal
}

func getAsk (currRow int) decimal.Decimal{

  currPrice := fileSlice[0][currRow][7] //Change this
	currPriceDecimal, err := decimal.NewFromString(currPrice)

	if err != nil {
		panic(err)
	}

	return currPriceDecimal
}
