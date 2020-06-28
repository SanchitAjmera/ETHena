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

func getInfo(ctx context.Context, lunoClient *luno.Client, tick luno.GetTickerRequest, timeInMins int) (open decimal.Decimal, close decimal.Decimal, high decimal.Decimal, low decimal.Decimal) {

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

func main() {

	lunoClient := luno.NewClient()
	lunoClient.SetAuth("gwnarwdxreyag", "ZrCzrPO3IdcMq7t69a5iPUl-JyDAGGxauF0HumJD34s")
	ctx := context.Background()
	tick := luno.GetTickerRequest{Pair: "XBTGBP"}

	getInfo(ctx, lunoClient, tick, 1)

}
