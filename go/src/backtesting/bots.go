package main

import (
        "github.com/luno/luno-go/decimal"
)


//£-3698.61 per day
func verySimpleBot(nextPrice decimal.Decimal, lastPrice *decimal.Decimal) int {
	returnVal := nextPrice.Sub(*lastPrice).Sign()
	*lastPrice = nextPrice
	return returnVal
}

func SMEBot(state *state_t) {
  action := checkSME(state)

}
