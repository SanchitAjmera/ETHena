package main

import (
        "github.com/luno/luno-go/decimal"
)

//Tickers take the current row return the current price
type ticker func(int) decimal.Decimal

// struct for state of trader
type state_t struct {
  funds int
  assets int
  inventory [][]int
  historicalData [][][]string
  currentDay int
  metrics struct  {
    tickerTime int
    dataCacheLength int
  }
}
