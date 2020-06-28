package main

import (
        "strconv"
)

func SME(historicalData[] string, newData float64) float64 {

  n := len(historicalData)
	var sum float64 = 0

  for i := 0; i < n; i++ {

    decimalData, err := strconv.ParseFloat(historicalData[i], 8)
    if err != nil{
    }

    sum += decimalData
  }
  return sum / float64(n)
}

func checkSME(state *state_t) (int, float64) {

  newData:= getNextPrice(state.currentDay, state.historicalData)
  historicalDataSlice := state.historicalData[sheetNum][state.currentDay - state.metrics.dataCacheLength:state.currentDay][priceCol]
  average := SME(historicalDataSlice, newData)

  if average > 0 {
    return 1, newData

  } else if average < 0 {
    return -1, newData
  }
    return 0, newData
}
