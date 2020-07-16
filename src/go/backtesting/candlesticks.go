package main

import (
	"context"
	"fmt"
	"log"
	"time"

	luno "github.com/luno/luno-go"
	"github.com/luno/luno-go/decimal"
)

type Info struct {
	timeStamp luno.Time
	openAsk   decimal.Decimal
	closeAsk  decimal.Decimal
	maxAsk    decimal.Decimal
	minAsk    decimal.Decimal
	openBid   decimal.Decimal
	closeBid  decimal.Decimal
	maxBid    decimal.Decimal
	minBid    decimal.Decimal
}

// Assigning Global Variables
var ctx context.Context
var lunoClient *luno.Client

func getInfo(tick luno.GetTickerRequest, timeInMins int) Info {

	timeInSeconds := timeInMins * 60
	// timeInSeconds := timeInMins
	maxAsk := decimal.Zero()
	minAsk := decimal.NewFromInt64(1844674407370955200)
	maxBid := decimal.Zero()
	minBid := decimal.NewFromInt64(1844674407370955200)
	openAsk := decimal.Zero()
	closeAsk := decimal.Zero()
	openBid := decimal.Zero()
	closeBid := decimal.Zero()
	var timeStamp luno.Time

	for i := 0; i < timeInSeconds; i++ {
		res, err := lunoClient.GetTicker(ctx, &tick)
		if err != nil {
			log.Fatal(err)

		}
		// fmt.Println("Count : ", i, "  Ask: ", res.Ask, "  Bid: ", res.Bid, "  Time Stamp: ", res.Timestamp)

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
			timeStamp = res.Timestamp
		}

		time.Sleep(time.Second)
	}

	// fmt.Println("Open Ask : ", openAsk, "   Highest Ask : ", maxAsk, "  Lowest Ask  : ", minAsk, "  Close Ask : ", closeAsk)
	// fmt.Println("Open Bid : ", openBid, "   Highest Bid : ", maxBid, "  Lowest Bid  : ", minBid, "  Close Bid : ", closeBid)

	info := Info{timeStamp, openAsk, closeAsk, maxAsk, minAsk, openBid, closeBid, maxBid, minBid}

	return info
}

func upDownOrNothing(tick luno.GetTickerRequest, timeInMins int) int {
	// Returns 1 if prediction buy
	// Returns 0 if nothing
	// Returns -1 if prediction sell

	stack := []Info{getInfo(tick, timeInMins), getInfo(tick, timeInMins), getInfo(tick, timeInMins)}

	for {

		b1Op := stack[0].openAsk
		b1Cl := stack[0].closeAsk
		b1Max := stack[0].maxAsk
		b1Min := stack[0].minAsk

		b2Op := stack[1].openAsk
		b2Cl := stack[1].closeAsk
		b2Max := stack[1].maxAsk
		b2Min := stack[1].minAsk

		b3Op := stack[2].openAsk
		b3Cl := stack[2].closeAsk
		b3Max := stack[2].maxAsk
		b3Min := stack[2].minAsk

		fmt.Println("\n------------------ Time Stamp of Bar 3 : ", stack[2].timeStamp)

		if b2Max.Cmp(b1Max) == 1 && b2Max.Cmp(b3Max) == 1 && b2Min.Cmp(b1Min) == 1 && b2Min.Cmp(b3Min) == 1 {
			if b1Cl.Cmp(b1Op) == 1 && b3Op.Cmp(b3Cl) == 1 && b2Cl.Cmp(b2Op) == 1 {
				fmt.Println("Predicting to sell")
			}
		} else if b2Max.Cmp(b1Max) == -1 && b2Max.Cmp(b3Max) == -1 && b2Min.Cmp(b1Min) == -1 && b2Min.Cmp(b3Min) == -1 {
			if b1Cl.Cmp(b1Op) == -1 && b3Op.Cmp(b3Cl) == -1 && b2Cl.Cmp(b2Op) == -1 {
				fmt.Println("Predicting to buy")
			}
		} else {
			fmt.Println("Predicting to do nada")
		}

		fmt.Println("Bar 1 opening : ", b1Op, "  Bar 1 closing : ", b1Cl, "  Bar 1 high : ", b1Max, " Bar 1 low : ", b1Min)
		fmt.Println("Bar 2 opening : ", b2Op, "  Bar 2 closing : ", b2Cl, "  Bar 2 high : ", b2Max, " Bar 2 low : ", b2Min)
		fmt.Println("Bar 3 opening : ", b3Op, "  Bar 3 closing : ", b3Cl, "  Bar 3 high : ", b3Max, " Bar 3 low : ", b3Min)
		fmt.Println("------------------")

		stack = append(stack[1:], getInfo(tick, timeInMins))
	}
}

func main() {

	lunoClient = luno.NewClient()
	lunoClient.SetAuth("gwnarwdxreyag", "ZrCzrPO3IdcMq7t69a5iPUl-JyDAGGxauF0HumJD34s")
	ctx = context.Background()
	tick := luno.GetTickerRequest{Pair: "XBTGBP"}

	upDownOrNothing(tick, 1)
}
