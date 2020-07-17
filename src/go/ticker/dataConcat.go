package main

import (
  //"fmt"
  "github.com/luno/luno-go/decimal"
  "github.com/360EntSecGroup-Skylar/excelize"
  "strconv"
  "fmt"
)

func main(){
  columns := []string{"A", "B","C","D","E","F","G","H","I","J","K","L","M","N"}
  f := excelize.NewFile()
  f.SetCellValue("Sheet1", "A1", "XBTGBP bid")
  f.SetCellValue("Sheet1", "B1", "XBTGBP ask")
  f.SetCellValue("Sheet1", "C1", "ETHXBT bid")
  f.SetCellValue("Sheet1", "D1", "ETHXBT ask")
  f.SetCellValue("Sheet1", "E1", "XRPXBT bid")
  f.SetCellValue("Sheet1", "F1", "XRPXBT ask")
  f.SetCellValue("Sheet1", "G1", "XRPZAR bid")
  f.SetCellValue("Sheet1", "H1", "XRPZAR ask")
  f.SetCellValue("Sheet1", "I1", "BCHXBT bid")
  f.SetCellValue("Sheet1", "J1", "BCHXBT ask")
  f.SetCellValue("Sheet1", "K1", "LTCXBT bid")
  f.SetCellValue("Sheet1", "L1", "LTCXBT ask")
  f.SetCellValue("Sheet1", "M1", "XBTZAR bid")
  f.SetCellValue("Sheet1", "N1", "XBTZAR ask")
  for file := 0 ; file < 22 ; file++ {
    g, err := excelize.OpenFile("data_8to9_July/tickerData" + strconv.Itoa(file) + ".xlsx")
    if err != nil {
        fmt.Println(err)
        return
    }
    for row := 2; row < 62; row++{
      for _, column := range columns {
        cell, err := g.GetCellValue("Sheet1", column + strconv.Itoa(row))
        cellDec,_ := decimal.NewFromString(cell)
        fmt.Println("file",file," row",row," column", column, " value", cellDec.String())
        if err != nil {
          fmt.Println(err)
          return
        }
        f.SetCellValue("Sheet1", column + strconv.Itoa(row + file * 60),cellDec)
      }
    }
  }
  if err := f.SaveAs("tickerDataDay2.xlsx"); err != nil {
    println(err.Error())
  }
}
