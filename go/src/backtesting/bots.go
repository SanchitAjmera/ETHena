package main

import (
  "fmt"
)


func buy(state *state_t, newPrice float64 ){
  buyableStock := state.funds / newPrice

  if buyableStock != float64(0) {
    item := []float64{buyableStock, newPrice}
    state.inventory = append(state.inventory, item)
    state.funds -= buyableStock * newPrice
  }

  fmt.Println(".")
  fmt.Println("                   Bought %d stock at %d", buyableStock, newPrice)

}

func sell(state *state_t, newPrice float64) {
  sold := float64(0)

  if len(state.inventory) == 0 {
    fmt.Println(".")
    fmt.Println("                 Sold %d stock at %d", sold, newPrice)
    return
  }

  for i := 0; i < len(state.inventory); i++ {
    sold += state.inventory[i][0]
    state.funds += state.inventory[i][0] * newPrice
  }

  state.inventory = [][]float64{}
  fmt.Println(".")
  fmt.Println("                   Sold %d stocks at %d", sold, newPrice)

}


//£-3698.61 per day
func verySimpleBot(nextPrice float64, lastPrice *float64) float64 {
	returnVal := nextPrice - *lastPrice
	*lastPrice = nextPrice
	return returnVal
}

func SMEBot(state *state_t) {
  action, newPrice := checkSME(state)
  if action > 0 {
    buy(state, newPrice)
  } else if action < 0{
    sell(state, newPrice)
  }

  state.currentDay += 1
}
