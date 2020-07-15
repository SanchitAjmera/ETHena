package main

import (
	// "fmt"
	"context"
	"fmt"
	luno "github.com/luno/luno-go"
	"github.com/luno/luno-go/decimal"
	"time"
)

func getTickerRequest() (*luno.Client, *luno.GetTickerRequest) {
	lunoClient := luno.NewClient()
	lunoClient.SetAuth("e354ths3rzks7", "PdhKmE-IGaIqn_Ckb3-pyz7nItCZiRIhk2cBF57meC4")

	return lunoClient, &luno.GetTickerRequest{Pair: pair}
}

func getCurrBid() decimal.Decimal {
	res, err := client.GetTicker(context.Background(), reqPointer)
	if err != nil {
		panic(err)
		// fmt.Println(err)
		// time.Sleep(time.Minute)
		// return getCurrBid()
	}
	return res.Bid
}

func getCurrAsk() decimal.Decimal {
	res, err := client.GetTicker(context.Background(), reqPointer)
	if err != nil {
		panic(err)
		// fmt.Println(err)
		// time.Sleep(time.Minute)
		// return getCurrAsk()
	}
	return res.Ask
}

func getAsset(currency string) decimal.Decimal {
	balancesReq := luno.GetBalancesRequest{}
	balances, err := client.GetBalances(context.Background(), &balancesReq)
	if err != nil {
		panic(err)
		// fmt.Println(err)
		// time.Sleep(time.Minute)
		// getAsset(currency)
	}

	for _, accBalance := range balances.Balance {
		if accBalance.Asset == currency {
			return accBalance.Balance
		}
	}

	panic("Cannot retrieve account balance")
}

func getTicker() (decimal.Decimal, decimal.Decimal) {
	res, err := client.GetTicker(context.Background(), reqPointer)
	if err != nil {
		fmt.Println(err)
		time.Sleep(time.Minute)
		return getTicker()
	}
	return res.Ask, res.Bid
}

func getAssets(currency1 string , currency2 string) (decimal.Decimal, decimal.Decimal) {
	balancesReq := luno.GetBalancesRequest{}
	balances, err := client.GetBalances(context.Background(), &balancesReq)
	if err != nil {
		panic(err)
	}
	var return1 decimal.Decimal
	var return2 decimal.Decimal
	for _, accBalance := range balances.Balance {
		if accBalance.Asset == currency1 {
			return1 = accBalance.Balance
		}
		if accBalance.Asset == currency2 {
			return2 = accBalance.Balance
		}
	}
	return return1, return2
}
