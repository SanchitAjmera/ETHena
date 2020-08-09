package main

import (
	"fmt"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/luno/luno-go/decimal"
)

type info struct {
	timeStampRow int
	openAsk      decimal.Decimal
	closeAsk     decimal.Decimal
	maxAsk       decimal.Decimal
	minAsk       decimal.Decimal
	openBid      decimal.Decimal
	closeBid     decimal.Decimal
	maxBid       decimal.Decimal
	minBid       decimal.Decimal
}

var timeInPeriodMins = 30

func getInfoFromSpread(startRow int) info {
	// file, err := excelize.OpenFile("../../ticker/data_7to8_July/tickerData09072020.xlsx")
	file, err := excelize.OpenFile("../../ticker/data_8to9_July/tickerDataDay2.xlsx")
	if err != nil {
		panic(err)
	}

	maxAsk := decimal.Zero()
	minAsk := decimal.NewFromInt64(1844674407370955200)
	maxBid := decimal.Zero()
	minBid := decimal.NewFromInt64(1844674407370955200)
	openAsk := decimal.Zero()
	closeAsk := decimal.Zero()
	openBid := decimal.Zero()
	closeBid := decimal.Zero()

	upperRowBound := startRow + timeInPeriodMins

	for i := startRow; i <= upperRowBound; i++ {
		locationBid := "E" + strconv.Itoa(i)
		locationAsk := "F" + strconv.Itoa(i)
		cellBidStr, err := file.GetCellValue("Sheet1", locationBid)
		if err != nil {
			panic(err)
		}

		cellBid, _ := decimal.NewFromString(cellBidStr)

		cellAskStr, err := file.GetCellValue("Sheet1", locationAsk)
		if err != nil {
			panic(err)
		}

		cellAsk, _ := decimal.NewFromString(cellAskStr)

		if maxAsk.Cmp(cellAsk) == -1 {
			maxAsk = cellAsk
		}

		if maxBid.Cmp(cellBid) == -1 {
			maxBid = cellBid
		}

		if cellAsk.Cmp(minAsk) == -1 {
			minAsk = cellAsk
		}

		if cellBid.Cmp(minBid) == -1 {
			minBid = cellBid
		}

		if i == startRow {
			openAsk = cellAsk
			openBid = cellBid
		}

		if i == upperRowBound {
			closeAsk = cellAsk
			closeBid = cellBid
		}
	}
	info := info{upperRowBound, openAsk, closeAsk, maxAsk, minAsk, openBid, closeBid, maxBid, minBid}
	// fmt.Println("Row: ", info.timeStampRow, " OpenAsk: ", info.openAsk, " CloseAsk: ", info.closeAsk, " MaxAsk: ", info.maxAsk, " MinAsk: ", info.minAsk)
	return info
}

func upDownOrNothing() {

	stack := []info{getInfoFromSpread(2), getInfoFromSpread(2 + 30), getInfoFromSpread(2 + 60)}

	for currentRow := 62; currentRow < 1322; currentRow += timeInPeriodMins {

		b1Op := stack[0].openAsk
		b1Cl := stack[0].closeAsk
		b1Max := stack[0].maxAsk
		b1Min := stack[0].minAsk

		b2Op := stack[1].openAsk
		b2Cl := stack[1].closeAsk
		b2Max := stack[1].maxAsk
		b2Min := stack[1].minAsk

		b3Op := stack[2].openAsk
		b3Cl := stack[2].closeAsk
		b3Max := stack[2].maxAsk
		b3Min := stack[2].minAsk

		fmt.Println("\n------------------ Time Stamp of Bar 3 : ", stack[2].timeStampRow)

		if b1Cl.Cmp(b1Op) == -1 && b2Min.Cmp(b1Min) == -1 && b2Min.Cmp(b3Min) == -1 && b3Cl.Cmp(b1Max) == 1 && b3Cl.Cmp(b2Max) == 1 {
			fmt.Println("Predicting to buy")
		} else if b1Cl.Cmp(b1Op) == 1 && b2Min.Cmp(b1Min) == 1 && b2Min.Cmp(b3Min) == 1 && b3Cl.Cmp(b1Max) == -1 && b3Cl.Cmp(b2Max) == -1 {
			fmt.Println("Predicting to sell")
		} else {
			fmt.Println("Predicting to do nada")
		}

		fmt.Println("Bar 1 opening : ", b1Op, "  Bar 1 closing : ", b1Cl, "  Bar 1 high : ", b1Max, " Bar 1 low : ", b1Min)
		fmt.Println("Bar 2 opening : ", b2Op, "  Bar 2 closing : ", b2Cl, "  Bar 2 high : ", b2Max, " Bar 2 low : ", b2Min)
		fmt.Println("Bar 3 opening : ", b3Op, "  Bar 3 closing : ", b3Cl, "  Bar 3 high : ", b3Max, " Bar 3 low : ", b3Min)
		fmt.Println("------------------")

		stack = append(stack[1:], getInfoFromSpread(currentRow))
	}
}

func main() {
	fmt.Println("Hello world")

	upDownOrNothing()

}
