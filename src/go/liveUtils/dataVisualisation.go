package liveUtils

import (
  "strconv"
)

var f *excelize.File

func ClosePrevFile(fileName string) {
  if f == nil {return}
  if err := f.SaveAs(fileName+".xlsx"); err != nil {
    panic(err)
  }
}

func SetUpNewFile() {
	f = excelize.NewFile()
	f.SetCellValue("Sheet1", "A1", "Sr No.")
	f.SetCellValue("Sheet1", "B1", "Ask")
	f.SetCellValue("Sheet1", "C1", "Bid")
	f.SetCellValue("Sheet1", "D1", "RSI")
	f.SetCellValue("Sheet1", "E1", "ReadyToBuy")
	f.SetCellValue("Sheet1", "F1", "StopLoss")
	f.SetCellValue("Sheet1", "G1", "Ready To Buy Price")
	f.SetCellValue("Sheet1", "H1", "Ready To Sell Price")
}

func PopulateFile(b *rsiBot, ask decimal.Decimal, bid decimal.Decimal, rsi decimal.Decimal){
	rowNum = b.NumOfDecisions + 1
	f.SetCellValue("Sheet1", "A"+strconv.FormatInt(rowNum, 10), b.NumOfDecisions)
	f.SetCellValue("Sheet1", "B"+strconv.FormatInt(rowNum, 10), ask)
	f.SetCellValue("Sheet1", "C"+strconv.FormatInt(rowNum, 10), bid)
	f.SetCellValue("Sheet1", "D"+strconv.FormatInt(rowNum, 10), rsi)
	f.SetCellValue("Sheet1", "E"+strconv.FormatInt(rowNum, 10), b.ReadyToBuy)
  if b.ReadyToBuy {
    f.SetCellValue("Sheet1", "G"+strconv.FormatInt(rowNum, 10), ask)
  } else {
    f.SetCellValue("Sheet1", "F"+strconv.FormatInt(rowNum, 10), b.StopLoss)
    f.SetCellValue("Sheet1", "H"+strconv.FormatInt(rowNum, 10), bid)
  }


}
