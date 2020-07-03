package main

import (
	"context"
	"fmt"
	"log"
	"time"

	luno "github.com/luno/luno-go"
	"github.com/luno/luno-go/decimal"
)

// Assigning Global Variables
var ctx context.Context
var lunoClient *luno.Client

func getInfo(tick luno.GetTickerRequest, timeInMins int) (decimal.Decimal, decimal.Decimal, decimal.Decimal, decimal.Decimal) {

	timeInSeconds := timeInMins * 60
	maxAsk := decimal.Zero()
	minAsk := decimal.NewFromInt64(1844674407370955200)
	maxBid := decimal.Zero()
	minBid := decimal.NewFromInt64(1844674407370955200)
	openAsk := decimal.Zero()
	closeAsk := decimal.Zero()
	openBid := decimal.Zero()
	closeBid := decimal.Zero()

	for i := 0; i < timeInSeconds; i++ {
		res, err := lunoClient.GetTicker(ctx, &tick)
		if err != nil {
			log.Fatal(err)

		}
		fmt.Println("Count : ", i, "  Ask: ", res.Ask, "  Bid: ", res.Bid, "  Time Stamp: ", res.Timestamp)

		if maxAsk.Cmp(res.Ask) == -1 {
			maxAsk = res.Ask
		}

		if maxBid.Cmp(res.Bid) == -1 {
			maxBid = res.Bid
		}

		if res.Ask.Cmp(minAsk) == -1 {
			minAsk = res.Ask
		}

		if res.Bid.Cmp(minBid) == -1 {
			minBid = res.Bid
		}

		if i == 0 {
			openAsk = res.Ask
			openBid = res.Bid
		}

		if i == 59 {
			closeAsk = res.Ask
			closeBid = res.Bid
		}

		time.Sleep(time.Second)
	}

	fmt.Println("Open Ask : ", openAsk, "   Highest Ask : ", maxAsk, "  Lowest Ask  : ", minAsk, "  Close Ask : ", closeAsk)
	fmt.Println("Open Bid : ", openBid, "   Highest Bid : ", maxBid, "  Lowest Bid  : ", minBid, "  Close Bid : ", closeBid)

	return openAsk, closeAsk, maxAsk, minAsk
}

func upDownOrNothing(tick luno.GetTickerRequest, timeInMins int) int {
	// Returns 1 if prediction buy
	// Returns 0 if nothing
	// Returns -1 if prediction sell
	b1Op, b1Cl, b1Max, b1Min := getInfo(tick, timeInMins)
	b2Op, b2Cl, b2Max, b2Min := getInfo(tick, timeInMins)
	b3Op, b3Cl, b3Max, b3Min := getInfo(tick, timeInMins)

	if b2Max.Cmp(b1Max) == 1 && b2Max.Cmp(b3Max) == 1 && b2Min.Cmp(b1Min) == 1 && b2Min.Cmp(b3Min) == 1 {
		if b1Cl.Cmp(b1Op) == 1 && b3Op.Cmp(b3Cl) == 1 && b2Cl.Cmp(b2Op) == 1 {
			return -1
		}
	} else if b2Max.Cmp(b1Max) == -1 && b2Max.Cmp(b3Max) == -1 && b2Min.Cmp(b1Min) == -1 && b2Min.Cmp(b3Min) == -1 {
		if b1Cl.Cmp(b1Op) == -1 && b3Op.Cmp(b3Cl) == -1 && b2Cl.Cmp(b2Op) == -1 {
			return 1
		}
	}
	return 0

}

func main() {

	lunoClient = luno.NewClient()
	lunoClient.SetAuth("gwnarwdxreyag", "ZrCzrPO3IdcMq7t69a5iPUl-JyDAGGxauF0HumJD34s")
	ctx = context.Background()
	tick := luno.GetTickerRequest{Pair: "XBTGBP"}

	fmt.Println(upDownOrNothing(tick, 1))

}
