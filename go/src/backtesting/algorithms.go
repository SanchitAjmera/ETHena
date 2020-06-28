package main

import (
        "strconv"
        "fmt"
)

func SME(state *state_t) (int, float64) {

  newData, timeStamp:= getNextPrice(state.currentDay, state.historicalData)

  if timeStamp == float64(0) {
    return 0, newData
  }
	var sum float64 = 0

  for i := 0; i < state.metrics.dataCacheLength; i++ {
    // TODO : should moving average include the current price
    decimalData, err := strconv.ParseFloat(state.historicalData[sheetNum][state.currentDay - i][priceCol], 8)
    if err != nil{
      panic(err)
    }

    sum += decimalData
  }

  average := sum / float64(state.metrics.dataCacheLength)


  if average + state.metrics.offset  <  newData {
    fmt.Println(".")
    fmt.Println("                   Price:          £"  , newData,
  							"\n                   Time:          " , timeStamp)
    fmt.Println("                   Moving Average: £",average)
    return -1, newData

  } else if average - state.metrics.offset > newData {
    fmt.Println(".")
    fmt.Println("                   Price:          £"  , newData,
  							"\n                   Time:          " , timeStamp)
    fmt.Println("                   Moving Average: £",average)
    return 1, newData
  }

  return 0, newData
}
