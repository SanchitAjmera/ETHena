package liveUtils

import (
	// "fmt"
	"context"
	"fmt"
	"time"

	luno "github.com/luno/luno-go"
	"github.com/luno/luno-go/decimal"
)

// Global Variables
var Client *luno.Client
var ReqPointer *luno.GetTickerRequest
var Pair string

func GetTickerRequest() (*luno.Client, *luno.GetTickerRequest) {
	lunoClient := luno.NewClient()
	lunoClient.SetAuth("g8gve3tacy6z5", "BL4y6PtBsWJHrqDZChS7bcEszLPWpt-z82Sls5S_a0g")

	return lunoClient, &luno.GetTickerRequest{Pair: Pair}
}

func getCurrBid() decimal.Decimal {
	res, err := Client.GetTicker(context.Background(), ReqPointer)
	if err != nil {
		panic(err)
		// fmt.Println(err)
		// time.Sleep(time.Minute)
		// return getCurrBid()
	}
	return res.Bid
}

func GetCurrAsk() decimal.Decimal {
	res, err := Client.GetTicker(context.Background(), ReqPointer)
	if err != nil {
		panic(err)
		// fmt.Println(err)
		// time.Sleep(time.Minute)
		// return GetCurrAsk()
	}
	return res.Ask
}

func getAsset(currency string) decimal.Decimal {
	balancesReq := luno.GetBalancesRequest{}
	balances, err := Client.GetBalances(context.Background(), &balancesReq)
	if err != nil {
		fmt.Println(err)
		time.Sleep(time.Minute)
		getAsset(currency)
	}

	for _, accBalance := range balances.Balance {
		if accBalance.Asset == currency {
			return accBalance.Balance
		}
	}

	panic("Cannot retrieve account balance")
}

func getTicker() (decimal.Decimal, decimal.Decimal) {
	res, err := Client.GetTicker(context.Background(), ReqPointer)
	if err != nil {
		fmt.Println(err)
		time.Sleep(time.Minute)
		return getTicker()
	}
	return res.Ask, res.Bid
}

func getAssets(currency1 string, currency2 string) (decimal.Decimal, decimal.Decimal) {
	balancesReq := luno.GetBalancesRequest{}
	balances, err := Client.GetBalances(context.Background(), &balancesReq)
	if err != nil {
		fmt.Println(err)
		time.Sleep(time.Minute)
		getAssets(currency1, currency2)
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
