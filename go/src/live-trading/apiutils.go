package main

import (
  "fmt"
  "context"
  luno "github.com/luno/luno-go"
	"github.com/luno/luno-go/decimal"
)

func getTickerRequest() (*luno.Client, *luno.GetTickerRequest){
  lunoClient := luno.NewClient()
  lunoClient.SetAuth("gc4t3v7f7tg5a", "OD8HmPz7QGtzqCZUc04JIfPKYUXLrkN1lzxppTc7cSs")

  return lunoClient, &luno.GetTickerRequest{Pair: pair}
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

func getAsset (currency string) decimal.Decimal{
	balancesReq := GetBalancesRequest{[]string{currency}}
	balances, err := client.GetBalances(context.Background(), &balancesReq)

	if err != nil {panic(err)}

	return balances[0].Balance
}
