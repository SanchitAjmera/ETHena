package liveUtils

import (
	"context"
	"log"
	"time"

	luno "github.com/luno/luno-go"
	"github.com/luno/luno-go/decimal"
)

// Global Variables
var Client *luno.Client
var PairName string

func CreateClient() *luno.Client {
	lunoClient := luno.NewClient()
	lunoClient.SetAuth("bgh3q5ss8yvj9", "Q9AqNYS4a4ke-4Yxw_JwtbUbih2mFjITm40luA6Y4M0")
	lunoClient.SetTimeout(2 * time.Minute)
	return lunoClient
}

func GetCurrAsk() decimal.Decimal {
	return getTickerRes().Ask
}

func getAsset(currency string) decimal.Decimal {
	balancesReq := luno.GetBalancesRequest{}
	balances, err := Client.GetBalances(context.Background(), &balancesReq)
	if err != nil {
		log.Println(err)
		time.Sleep(2 * time.Second)
		return getAsset(currency)
	}

	for _, accBalance := range balances.Balance {
		if accBalance.Asset == currency {
			return accBalance.Balance
		}
	}

	panic("Cannot retrieve account balance")
}

func getTickerRes() luno.GetTickerResponse {
	reqPointer := luno.GetTickerRequest{Pair: PairName}
	res, err := Client.GetTicker(context.Background(), &reqPointer)
	if err != nil {
		log.Println(err)
		time.Sleep(2 * time.Second)
		return getTickerRes()
	}
	return *res
}

func getAssets(currency1 string, currency2 string) (decimal.Decimal, decimal.Decimal) {
	balancesReq := luno.GetBalancesRequest{}
	balances, err := Client.GetBalances(context.Background(), &balancesReq)
	if err != nil {
		log.Println(err)
		time.Sleep(2 * time.Second)
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
