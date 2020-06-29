package main

import ("fmt"
        "github.com/tealeg/xlsx"
      )

//var priceCol int
//var timeCol int
//var sheetNum int


func main () {

  /*initialising global variables
  priceCol = 7
  timeCol = 0
  sheetNum = 0
*/
//	var getNextPrice ticker = tickerFunc()
  fileSlice, err := xlsx.FileToSlice("recentAPIdata.xlsx")

  if err != nil {
    panic(err)
  }

  //tester(30,100,2000,fileSlice)


  for duration := 1000;  duration< 16000; duration+= 1000 {
    maxFunds := float64(0)
    maxMinutes := float64(0)
    maxOffset := float64(0)
    // minutes is the length which the EMA/SME is calculated over
    // and how often the bot queries for price changes
    for minutes := 1; minutes < 120; minutes+=1 {
      for offset := 5;  offset< 200; offset+=5 {
      //  fmt.Println(".")
      //  fmt.Println("               ", minutes, offset, duration)
        result := tester(minutes, offset, duration,fileSlice)
        if result > maxFunds  {
          maxFunds = result
          maxMinutes = float64(minutes)
          maxOffset = float64(offset)
        }
      }
    }

    fmt.Println("  Max Funds £",maxFunds," over duration:", duration," with offset",maxOffset, " and minute intervals of", maxMinutes)
  }
}
