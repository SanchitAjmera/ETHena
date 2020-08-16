package main

import (
	"fmt"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/luno/luno-go/decimal"
)

var timeInPeriodMins = 30

func getInfoFromSpread(startRow int) candlestick {
	// file, err := excelize.OpenFile("../ticker/data_7to8_July/tickerData09072020.xlsx")
	file, err := excelize.OpenFile("../../ticker/data/18to19-July.xlsx")
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
	info := candlestick{openAsk, closeAsk, maxAsk, minAsk, openBid, closeBid, maxBid, minBid}
	return info
}

func rev123(stick1 candlestick, stick2 candlestick, stick3 candlestick) bool {
	b1Op := stick1.openAsk
	b1Cl := stick1.closeAsk
	b1Max := stick1.maxAsk
	b1Min := stick1.minAsk

	// b2Op := stick2.openAsk
	// b2Cl := stick2.closeAsk
	b2Min := stick2.minAsk
	b2Max := stick2.maxAsk

	// b3Op := stick3.openAsk
	b3Cl := stick3.closeAsk
	// b3Max := stick3.maxAsk
	b3Min := stick3.minAsk

	//For buying
	return b1Cl.Cmp(b1Op) == -1 && b2Min.Cmp(b1Min) == -1 && b2Min.Cmp(b3Min) == -1 && b3Cl.Cmp(b1Max) == 1 && b3Cl.Cmp(b2Max) == 1
}

func hammer(stick candlestick) bool {
	op := stick.openAsk
	cl := stick.closeAsk
	// max := stick.maxAsk
	min := stick.minAsk

	diffClOp := cl.Sub(op)
	hammerScale, _ := decimal.NewFromString("2")
	diffOpMin := (op.Sub(min)).Mul(hammerScale)

	return op.Cmp(cl) == -1 && diffOpMin.Cmp(diffClOp) == 1
}

func inverseHammer(stick candlestick) bool {
	op := stick.openAsk
	cl := stick.closeAsk
	max := stick.maxAsk
	// min := stick.minAsk

	diffClOp := cl.Sub(op)
	hammerScale, _ := decimal.NewFromString("2")
	diffMaxCl := (max.Sub(cl)).Mul(hammerScale)

	return op.Cmp(cl) == -1 && (diffMaxCl).Cmp(diffClOp) == 1
}

func whiteSlaves(stick1 candlestick, stick2 candlestick, stick3 candlestick) bool {
	return stick1.openAsk.Cmp(stick1.closeAsk) == -1 && stick2.openAsk.Cmp(stick2.closeAsk) == -1 && stick3.openAsk.Cmp(stick3.closeAsk) == -1
}

func morningStar(stick1 candlestick, stick2 candlestick, stick3 candlestick) bool {
	b1Op := stick1.openAsk
	b1Cl := stick1.closeAsk

	b2Op := stick2.openAsk
	b2Cl := stick2.closeAsk

	b3Op := stick3.openAsk
	b3Cl := stick3.closeAsk

	diffb1 := b1Op.Sub(b1Cl)
	diffb3 := b3Cl.Sub(b3Op)
	scale, _ := decimal.NewFromString("3")
	diffb2 := (b2Op.Sub(b2Cl)).Mul(scale)

	return b1Op.Cmp(b1Cl) == 1 && b2Op.Cmp(b2Cl) == 1 && b3Op.Cmp(b3Cl) == -1 && diffb1.Cmp(diffb2) == 1 && diffb3.Cmp(diffb2) == 1
}

func tradeBacktester() {

	stack := []candlestick{getInfoFromSpread(2), getInfoFromSpread(2 + 30), getInfoFromSpread(2 + 60)}

	for currentRow := 62; currentRow < 1322; currentRow += timeInPeriodMins {

		fmt.Println("123Rev : ", rev123(stack[0], stack[1], stack[2]))
		fmt.Println("Hammer : ", hammer(stack[2]))
		fmt.Println("Inverse Hammer : ", inverseHammer(stack[2]))
		fmt.Println("White Slaves : ", whiteSlaves(stack[0], stack[1], stack[2]))
		fmt.Println("Morningstar : ", morningStar(stack[0], stack[1], stack[2]))

		stack = append(stack[1:], getInfoFromSpread(currentRow))
	}
}

// func main() {
// 	trade()
// }
