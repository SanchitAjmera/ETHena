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


  for duration := 1000;  duration< 31000; duration+= 1000 {
    maxFunds := float64(0)
    maxMinutes := float64(0)
    maxOffset := float64(0)

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

  // MAX seems to be funds: 105397.6732124628  minutes: 330   offset: 85
  //106683.81605602588 26 20

}
