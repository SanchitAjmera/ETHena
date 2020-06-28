package main

func SME(historicalData []int, newData int) int {
  n := len(historicalData)
  sum := 0
  for i := 0; i<n; i++ {
    sum += historicalData[i]
  }
  return sum/n
}

func checkSME(*state state_t){

  average := SME(historicalData, newData)

  if average > 0 {
    return 1

  } else if average < 0 {

    return -1
  }

  return 0
}
