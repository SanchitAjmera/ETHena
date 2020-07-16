package main

import (
  "fmt"
  "context"
  luno "github.com/luno/luno-go"
	"github.com/luno/luno-go/decimal"
)

func getTickerRequest(pairName string) (*luno.Client, *luno.GetTickerRequest){
  lunoClient := luno.NewClient()
  lunoClient.SetAuth("gwnarwdxreyag", "ZrCzrPO3IdcMq7t69a5iPUl-JyDAGGxauF0HumJD34s")

  return lunoClient, &luno.GetTickerRequest{Pair: pairName}
}

func getCurrBid() decimal.Decimal{
    res, err := client.GetTicker(context.Background(), reqPointer)
    if err != nil {
      fmt.Errorf("Error in retrieving bid")
    }
    return res.Bid
}

func getCurrAsk() decimal.Decimal{
    res, err := client.GetTicker(context.Background(), reqPointer)
    if err != nil {
      fmt.Errorf("Error in retrieving ask")
    }
    return res.Ask
}
