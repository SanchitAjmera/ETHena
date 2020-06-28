package main

import (
        "github.com/luno/luno-go/decimal"
)

//Tickers take the current row return the current price
type ticker func(int) decimal.Decimal

// struct for state of trader
type state_t struct {
  funds float64
  assets float64
  inventory [][]float64
  historicalData [][][]string
  currentDay int
  metrics struct  {
    offset float64
    tickerTime int
    dataCacheLength int
  }
}
