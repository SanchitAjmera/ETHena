package main

import (
        "github.com/luno/luno-go/decimal"
)

//Tickers take the current row return the current price
type ticker func(int) decimal.Decimal
