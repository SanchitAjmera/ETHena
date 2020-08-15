package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/luno/luno-go"
	"github.com/luno/luno-go/decimal"
)

// Client is the luno client
var Client *luno.Client

// PairName is pair
var PairName string

type candleBot struct {
	TradingPeriod int64           // How often the bot calculates a result
	TradesMade    int64           // total number of trades executed
	ReadyToBuy    bool            // holds the state of the bot
	BuyPrice      decimal.Decimal // price we bought at
	PrevOrder     string
	StopLoss      decimal.Decimal
	StopLossMult  decimal.Decimal
}

type candlestick struct {
	openAsk  decimal.Decimal
	closeAsk decimal.Decimal
	maxAsk   decimal.Decimal
	minAsk   decimal.Decimal
	openBid  decimal.Decimal
	closeBid decimal.Decimal
	maxBid   decimal.Decimal
	minBid   decimal.Decimal
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

// function to execute buying of items
func buy(b *candleBot, currAsk decimal.Decimal) {
	time.Sleep(time.Second * 2)

	_, startFunds := getAssets("ETH", "XBT")
	price := currAsk.Sub(decimal.NewFromFloat64(0.000001, 8))
	buyableStock := startFunds.Div(price, 8)
	buyableStock = buyableStock.ToScale(2)

	// checking if there are no funds available
	if buyableStock.Sign() == 0 {
		log.Println("Not enough funds available")
		b.ReadyToBuy = false
		return
	}
	//Create limit order
	req := luno.PostLimitOrderRequest{
		Pair:   PairName,
		Price:  price,
		Type:   "BID", //We are putting in a bid to buy at the ask price
		Volume: buyableStock,
		//BaseAccountId: --> Not needed until using multiple strategies
		//CounterAccountId: --> Same as above
		PostOnly: true,
	}
	res, err := Client.PostLimitOrder(context.Background(), &req)
	for err != nil {
		log.Println(err)
		time.Sleep(time.Second * 30)
		res, err = Client.PostLimitOrder(context.Background(), &req)
	}
	log.Println("BUY - order ", res.OrderId, " placed at ", price)
	b.PrevOrder = res.OrderId
	b.ReadyToBuy = false
	b.TradesMade++
	b.StopLoss = price
	b.BuyPrice = price
	// wait till order has gone through
	log.Println("Waiting for buy order to be partially filled")
	counter := 0
	for {
		time.Sleep(2 * time.Second)
		counter++
		if getOrderStatus(b.PrevOrder) == "COMPLETE" {
			return
		}
		if counter > 15 {
			b.TradesMade--
			log.Println("Buy order timed out. Retrying")
			buy(b, currAsk)
			return
		}
	}
}

func sell(b *candleBot, currBid decimal.Decimal) {
	time.Sleep(time.Second * 2)

	startStock, _ := getAssets("ETH", "XBT")
	startStock = startStock.ToScale(2)

	// checking if there are no stock available
	log.Println("startstock after scale: ", startStock)
	if startStock.Sign() == 0 {
		log.Println("Not enough stock available")
		b.ReadyToBuy = true
		return
	}
	price := currBid.Add(decimal.NewFromFloat64(0.000001, 8))

	req := luno.PostLimitOrderRequest{
		Pair:   PairName,
		Price:  price,
		Type:   "ASK",
		Volume: startStock,
		//BaseAccountId: --> Not needed until using multiple strategies
		//CounterAccoundId: --> Same as above
		PostOnly: true,
	}
	res, err := Client.PostLimitOrder(context.Background(), &req)
	for err != nil {
		log.Println(err)
		time.Sleep(2 * time.Second)
		res, err = Client.PostLimitOrder(context.Background(), &req)
	}

	log.Println("SELL - order ", res.OrderId, " placed at ", price)
	b.PrevOrder = res.OrderId
	b.ReadyToBuy = true
	b.TradesMade++
	log.Println("Waiting for sell order to be partially filled")
	counter := 0
	for {
		time.Sleep(2 * time.Second)
		counter++
		if getOrderStatus(b.PrevOrder) == "COMPLETE" {
			return
		}
		if counter > 15 {
			b.TradesMade--
			log.Println("Sell order timed out. Retrying")
			sell(b, currBid)
			return
		}
	}
}

func getOrderStatus(id string) luno.OrderState {
	orderReq := luno.GetOrderRequest{
		Id: id,
	}

	res, err := Client.GetOrder(context.Background(), &orderReq)
	if err != nil {
		log.Println(err)
		time.Sleep(2 * time.Second)
		return getOrderStatus(id)
	}
	return res.State
}

func getCandleStick(b *candleBot) candlestick {
	maxAsk := decimal.Zero()
	minAsk := decimal.NewFromInt64(1844674407370955200)
	openAsk := decimal.Zero()
	closeAsk := decimal.Zero()
	maxBid := decimal.Zero()
	minBid := decimal.NewFromInt64(1844674407370955200)
	openBid := decimal.Zero()
	closeBid := decimal.Zero()

	for i := 0; int64(i) <= b.TradingPeriod; i++ {
		res := getTickerRes()
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

		if int64(i) == b.TradingPeriod {
			closeAsk = currAsk
			closeBid = currBid
		}
		time.Sleep(time.Second)
	}
	stick := candlestick{openAsk, closeAsk, maxAsk, minAsk, openBid, closeBid, maxBid, minBid}
	return stick
}

func tradeLive(b *candleBot) {
	fmt.Println("Processing 1")
	initStick1 := getCandleStick(b)
	fmt.Println("Processing 2")

	initStick2 := getCandleStick(b)
	fmt.Println("Processing 3")
	initStick3 := getCandleStick(b)

	stack := []candlestick{initStick1, initStick2, initStick3}

	for {

		if b.ReadyToBuy {
			fmt.Println("123Rev : ", rev123(stack[0], stack[1], stack[2]))
			fmt.Println("Hammer : ", hammer(stack[2]))
			fmt.Println("Inverse Hammer : ", inverseHammer(stack[2]))
			fmt.Println("White Slaves : ", whiteSlaves(stack[0], stack[1], stack[2]))
			fmt.Println("Morningstar : ", morningStar(stack[0], stack[1], stack[2]))
			if rev123(stack[0], stack[1], stack[2]) || hammer(stack[2]) || inverseHammer(stack[2]) || whiteSlaves(stack[0], stack[1], stack[2]) || morningStar(stack[0], stack[1], stack[2]) {
				buy(b, stack[2].closeAsk)
			}
		} else {
			bound := stack[2].closeBid.Mul(b.StopLossMult)
			currBid := stack[2].closeBid
			if (currBid.Cmp(b.BuyPrice) == 1 && currBid.Cmp(b.StopLoss) == -1) || currBid.Cmp(b.BuyPrice.Mul(decimal.NewFromFloat64(0.99, 8))) == -1 {
				sell(b, currBid)
			} else if bound.Cmp(b.StopLoss) == 1 {
				b.StopLoss = bound
				log.Println("Stoploss changed to: ", b.StopLoss)
			}
		}
		// fmt.Println(stack)
		stack = append(stack[1:], getCandleStick(b))
	}
}

func main() {

	PairName = "ETHXBT"
	// live.InitialiseKeys()
	// live.User = strings.ToUpper("devam")
	// live.Client = live.CreateClient()

	Client = luno.NewClient()
	Client.SetAuth("mggh7nx5v5vzn",
		"DiHrN4Lqu27eCajdCTBEKU4H-oIFAFR4_k1eRlx5Kho")

	bot := candleBot{
		TradingPeriod: 60,
		TradesMade:    0,
		ReadyToBuy:    true,
		BuyPrice:      decimal.Zero(),
		StopLoss:      decimal.Zero(),
		StopLossMult:  decimal.NewFromFloat64(0.9975, 8),
	}

	tradeLive(&bot)
	// tradeBacktester()
}
