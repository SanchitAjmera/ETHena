package main

import (
        "github.com/luno/luno-go/decimal"
)

//Tickers take the current row return the current price
type ticker func(int) decimal.Decimal

// global variable for data
type state_t struct {
  funds int
  assets int
  inventory [][]int
  historicalData [][][]string
  currentDay int
  struct metrics {
    tickerTime int
  }
}
