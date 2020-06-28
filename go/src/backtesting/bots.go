package main

import (
)


//£-3698.61 per day
func verySimpleBot(nextPrice float64, lastPrice *float64) float64 {
	returnVal := nextPrice - *lastPrice
	*lastPrice = nextPrice
	return returnVal
}

func SMEBot(state *state_t) {
//  action := checkSME(state)

}
