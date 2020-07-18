package liveUtils

import (
	// "fmt"
	"context"
	"fmt"
	luno "github.com/luno/luno-go"
	"github.com/luno/luno-go/decimal"
	"time"
)

// Global Variables
var Client *luno.Client
var ReqPointer *luno.GetTickerRequest
var Pair string

func GetTickerRequest() (*luno.Client, *luno.GetTickerRequest) {
	lunoClient := luno.NewClient()
	lunoClient.SetAuth("e354ths3rzks7", "PdhKmE-IGaIqn_Ckb3-pyz7nItCZiRIhk2cBF57meC4")

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
	res, err := Client.GetTicker(context.Background(), ReqPointer)
	if err != nil {
		fmt.Println(err)
		time.Sleep(time.Minute)
		return getTicker()
	}
	return res.Ask, res.Bid
}

func getAssets(currency1 string , currency2 string) (decimal.Decimal, decimal.Decimal) {
	balancesReq := luno.GetBalancesRequest{}
	balances, err := Client.GetBalances(context.Background(), &balancesReq)
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
