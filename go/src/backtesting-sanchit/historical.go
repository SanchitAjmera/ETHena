package main

import (
	"fmt"
	"github.com/luno/luno-go/decimal"
	"github.com/tealeg/xlsx"
)

/*
TODO:
- change colNum far ask and bid
*/

//global variable for historicalData
var historicalData [][][]string

// function to process the csv file and return a 3d array of strings
// historicalData is of the form: [sheetNum][rowNum][colNum]
func parseXlsx() {
	fileSlice, err := xlsx.FileToSlice("recentAPIdata.xlsx")
	if err != nil {
		panic(err)
	}
	historicalData = fileSlice
}

// function to get the bid price from a given row in the excel spreadsheet
func getBid(currRow int64) decimal.Decimal {
	currPrice := historicalData[0][int(currRow)][7] //Change this
	// if data is non applicable skip this row
	if (currPrice == "NaN") {
		return getBid(currRow - 1)
	}

	currPriceDecimal, err := decimal.NewFromString(currPrice)

	if err != nil {
		panic(err)
	}
	return currPriceDecimal
}

// function to get the ask price from a given row in the excel spreadsheet
func getAsk(currRow int64) decimal.Decimal {
	currPrice := historicalData[0][currRow][7] //Change this
	// if data is non applicable skip this row
	if (currPrice == "NaN") {
		return getAsk(currRow - 1)
	}

	currPriceDecimal, err := decimal.NewFromString(currPrice)

	if err != nil {
		panic(err)
	}

	return currPriceDecimal
}

func printPortFolio(pf *portfolio) {
	fmt.Println("trade NO. :   ", pf.tradesMade)
	fmt.Println("funds : 			£", pf.funds)
	fmt.Println("stock : 			£", pf.stock)
	fmt.Println(".")
}
