package Utils

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
var User string
var ApiKeys map[string]([]string)

func CreateClient() *luno.Client {
	lunoClient := luno.NewClient()
	keyID := ApiKeys[User][0]
	key := ApiKeys[User][1]
	lunoClient.SetAuth(keyID, key)
	lunoClient.SetTimeout(2 * time.Minute)
	return lunoClient
}

func GetCurrAsk() decimal.Decimal {
	return GetTickerRes().Ask
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

func GetTickerRes() luno.GetTickerResponse {
	reqPointer := luno.GetTickerRequest{Pair: PairName}
	res, err := Client.GetTicker(context.Background(), &reqPointer)
	if err != nil {
		log.Println(err)
		time.Sleep(2 * time.Second)
		return GetTickerRes()
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

func GetCandleStick(tradinginterval int64) Candlestick {
	maxAsk := decimal.Zero()
	minAsk := decimal.NewFromInt64(1844674407370955200)
	openAsk := decimal.Zero()
	closeAsk := decimal.Zero()
	maxBid := decimal.Zero()
	minBid := decimal.NewFromInt64(1844674407370955200)
	openBid := decimal.Zero()
	closeBid := decimal.Zero()

	for i := 0; int64(i) <= tradinginterval; i++ {
		res := GetTickerRes()
		currAsk, currBid := res.Ask, res.Bid

		if maxAsk.Cmp(currAsk) == -1 {
			maxAsk = currAsk
		}

		if maxBid.Cmp(currBid) == -1 {
			maxBid = currBid
		}

		if currAsk.Cmp(minAsk) == -1 {
			minAsk = currAsk
		}

		if currBid.Cmp(minBid) == -1 {
			minBid = currBid
		}

		if i == 0 {
			openAsk = currAsk
			openBid = currBid
		}

		if int64(i) == tradinginterval {
			closeAsk = currAsk
			closeBid = currBid
		}
		time.Sleep(time.Second)
	}
	stick := Candlestick{openAsk, closeAsk, maxAsk, minAsk, openBid, closeBid, maxBid, minBid}
	return stick
}
