package main

import (
        "strconv"
        //"fmt"
)

func SME(state *state_t) (int, float64) {

  newData, timeStamp:= getNextPrice(state.currentDay, state.historicalData)

  if timeStamp == float64(0) {
    return 0, newData
  }
	var sum float64 = 0

  for i := 0; i < state.metrics.dataCacheLength; i++ {
    // TODO : should moving average include the current price
    decimalData, err := strconv.ParseFloat(state.historicalData[0][state.currentDay - i][7], 8)
    if err != nil{
      panic(err)
    }

    sum += decimalData
  }

  average := sum / float64(state.metrics.dataCacheLength)


  if average + state.metrics.offset  <  newData {
    /*if printing {
    fmt.Println(".")
    fmt.Println("                   Price:          £"  , newData,
      "\n                   Time:          " , timeStamp)
    fmt.Println("                   Moving Average: £",average)
    }*/

    return -1, newData

  } else if average - state.metrics.offset > newData {
    /*if printing {
    fmt.Println(".")
    fmt.Println("                   Price:          £"  , newData,
      "\n                   Time:          " , timeStamp)
    fmt.Println("                   Moving Average: £",average)
    }*/

    return 1, newData
  }

  return 0, newData
}



// the exponential moving average function
// always pass in a 0 value for pointer
func EMA(state *state_t, pointer int) float64 {

  if pointer == state.metrics.dataCacheLength {
    newPrice, _ := getNextPrice(state.currentDay - pointer+1, state.historicalData)
    return newPrice
  }

  newData, timeStamp := getNextPrice(state.currentDay - pointer+1, state.historicalData)
  pointer+=1

  if timeStamp == float64(0) {
    return EMA(state, pointer)
  }

  smoothing := float64(2)
  factor := smoothing/float64(state.metrics.dataCacheLength +1)

  return (newData * (factor)) + (EMA(state, pointer) * (1-factor))
}



func checkEMA(state *state_t) (int, float64){
  ema := EMA(state, 0)
  currentPrice, _ := getNextPrice(state.currentDay, state.historicalData)

  if ema + state.metrics.offset < currentPrice {
    //fmt.Println("                   EMA: £", ema)
    //fmt.Println("                   Price:          £"  , currentPrice)
    return 1, currentPrice
  } else if ema - state.metrics.offset > currentPrice{
    return -1, currentPrice
    //fmt.Println("                   EMA: £", ema)
    //fmt.Println("                   Price:          £"  , currentPrice)
  }
  return 0, currentPrice
}
