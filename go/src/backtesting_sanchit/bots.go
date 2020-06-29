package main

import (
//  "fmt"
)


func buy(state *state_t, newPrice float64, printing bool ){
  buyableStock := state.funds / newPrice

  if buyableStock != float64(0) {
    state.inventory[buyableStock] = newPrice
    state.funds -= buyableStock * newPrice
  }

  //if printing {
  //fmt.Println("                   Bought ", buyableStock, " stocks at ",newPrice)
  //}

}

func sell(state *state_t, newPrice float64, printing bool) {
  sold := float64(0)

  if len(state.inventory) == 0 {
    //if printing {
    //fmt.Println("                   Sold ", sold, " stocks at ",newPrice)
    //}
    return
  }

  for stock, price := range state.inventory {
    if price <= newPrice{
      sold += stock
      state.funds += stock * newPrice
      delete(state.inventory, stock)
    }
  }

  //if printing {
  //fmt.Println("                   Sold ", sold, " stocks at ",newPrice)
  //}

}


//£-3698.61 per day
func verySimpleBot(nextPrice float64, lastPrice *float64) float64 {
	returnVal := nextPrice - *lastPrice
	*lastPrice = nextPrice
	return returnVal
}



func SMEBot(state *state_t, printing bool) {
  action, newPrice := SME(state, printing)

  if action == 0 {
    state.currentDay += state.metrics.dataCacheLength
    return
  }

  //if printing {
//  fmt.Println("                   inventory:     ",state.inventory)
//  fmt.Println("                   funds:          £",state.funds)
  //}

  if action > 0 {
    buy(state, newPrice, printing)
  } else if action < 0 {
    sell(state, newPrice, printing)
  }

  state.currentDay += state.metrics.dataCacheLength
}
