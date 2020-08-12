package Utils

import (
	"github.com/luno/luno-go/decimal"
)

type Candlestick struct {
	OpenAsk  decimal.Decimal
	CloseAsk decimal.Decimal
	MaxAsk   decimal.Decimal
	MinAsk   decimal.Decimal
	OpenBid  decimal.Decimal
	CloseBid decimal.Decimal
	MaxBid   decimal.Decimal
	MinBid   decimal.Decimal
}

//struct for the rsiBot
type RsiBot struct {
	TradesMade           int64             // total number of trades executed
	NumOfDecisions       int64             // number of times the bot calculates
	StopLoss             decimal.Decimal   // variable stop loss
	StopLossMult         decimal.Decimal   // multiplier for stop loss
	OverSold             decimal.Decimal   // bound to tell the bot when to buy
	ReadyToBuy           bool              // false means ready to sell
	BuyPrice             decimal.Decimal   // stores most recent price we bought at
	UpEma                decimal.Decimal   // exponentially smoothed Wilder's MMA for upward change
	DownEma              decimal.Decimal   // exponentially smoothed Wilder's MMA for downward change
	PrevAsk              decimal.Decimal   // the previous recorded ask price
	PrevOrder            string            // stores order ID of most recent order
	RSITradingPeriod     int64             // No of past asks used to calculate RSI
	MACDTradingPeriodLR  int64             // Long term MACD period
	MACDTradingPeriodSR  int64             // Short term MACD period
	CandleTradingPeriod  int64             //candlestick period
	LongestTradingPeriod int64             // Longest of the trading periods
	MACDlongperiodavg    decimal.Decimal   // MACD long period average
	PastAsks             []decimal.Decimal // array of previous ask prices
	MACDshortperiodavg   decimal.Decimal   // MACD short period average
	TimeInterval         int64             //time between each test
	Stack                []Candlestick     //stack of candlesticks
	Offset               decimal.Decimal   //for offset bot
	OffsetTraingPeriod   int64             // trading period for offset bot
	BotString            string            //string to determine which bots are being usde
}

// function to calculate the Relative Strength Index
func GetRsi(PrevAsk decimal.Decimal, currAsk decimal.Decimal, UpEma decimal.Decimal, DownEma decimal.Decimal, period int64) (decimal.Decimal, decimal.Decimal, decimal.Decimal) {
	//iterating through elements of array and populating priceUp/Down arrays
	var upDiff decimal.Decimal
	var downDiff decimal.Decimal

	if currAsk.Cmp(PrevAsk) == 1 {
		//item is over sold
		upDiff = currAsk.Sub(PrevAsk)
		downDiff = decimal.Zero()
	} else if currAsk.Cmp(PrevAsk) == -1 {
		//item is over bought
		upDiff = decimal.Zero()
		downDiff = PrevAsk.Sub(currAsk)
	} else {
		upDiff = decimal.Zero()
		downDiff = decimal.Zero()
	}

	priceUp := Ema(UpEma, upDiff, period)
	priceDown := Ema(DownEma, downDiff, period)

	// check to see if average fall price is zero to avoid div by zero error
	if priceDown.Sign() == 0 {
		if priceUp.Sign() == 0 {
			return decimal.NewFromInt64(50), priceUp, priceDown
		} else {
			return decimal.NewFromInt64(100), priceUp, priceDown
		}
	}

	rs := priceUp.Div(priceDown, 16)
	rsiDen := rs.Add(decimal.NewFromInt64(1))
	// calculating rsi
	rsi := decimal.NewFromInt64(100).Sub(decimal.NewFromInt64(100).Div(rsiDen, 16))
	return rsi, priceUp, priceDown
}

func Rev123(stick1 Candlestick, stick2 Candlestick, stick3 Candlestick) bool {
	b1Op := stick1.OpenAsk
	b1Cl := stick1.CloseAsk
	b1Max := stick1.MaxAsk
	b1Min := stick1.MinAsk

	// b2Op := stick2.OpenAsk
	// b2Cl := stick2.CloseAsk
	b2Min := stick2.MinAsk
	b2Max := stick2.MaxAsk

	// b3Op := stick3.OpenAsk
	b3Cl := stick3.CloseAsk
	// b3Max := stick3.MaxAsk
	b3Min := stick3.MinAsk

	//For buying
	return b1Cl.Cmp(b1Op) == -1 && b2Min.Cmp(b1Min) == -1 && b2Min.Cmp(b3Min) == -1 && b3Cl.Cmp(b1Max) == 1 && b3Cl.Cmp(b2Max) == 1
}

func Hammer(stick Candlestick) bool {
	op := stick.OpenAsk
	cl := stick.CloseAsk
	// max := stick.MaxAsk
	min := stick.MinAsk

	diffClOp := cl.Sub(op)
	hammerScale, _ := decimal.NewFromString("2")
	diffOpMin := (op.Sub(min)).Mul(hammerScale)

	return op.Cmp(cl) == -1 && diffOpMin.Cmp(diffClOp) == 1
}

func InverseHammer(stick Candlestick) bool {
	op := stick.OpenAsk
	cl := stick.CloseAsk
	max := stick.MaxAsk
	// min := stick.MinAsk

	diffClOp := cl.Sub(op)
	hammerScale, _ := decimal.NewFromString("2")
	diffMaxCl := (max.Sub(cl)).Mul(hammerScale)

	return op.Cmp(cl) == -1 && (diffMaxCl).Cmp(diffClOp) == 1
}

func WhiteSlaves(stick1 Candlestick, stick2 Candlestick, stick3 Candlestick) bool {
	return stick1.OpenAsk.Cmp(stick1.CloseAsk) == -1 && stick2.OpenAsk.Cmp(stick2.CloseAsk) == -1 && stick3.OpenAsk.Cmp(stick3.CloseAsk) == -1
}

func MorningStar(stick1 Candlestick, stick2 Candlestick, stick3 Candlestick) bool {
	b1Op := stick1.OpenAsk
	b1Cl := stick1.CloseAsk

	b2Op := stick2.OpenAsk
	b2Cl := stick2.CloseAsk

	b3Op := stick3.OpenAsk
	b3Cl := stick3.CloseAsk

	diffb1 := b1Op.Sub(b1Cl)
	diffb3 := b3Cl.Sub(b3Op)
	scale, _ := decimal.NewFromString("3")
	diffb2 := (b2Op.Sub(b2Cl)).Mul(scale)

	return b1Op.Cmp(b1Cl) == 1 && b2Op.Cmp(b2Cl) == 1 && b3Op.Cmp(b3Cl) == -1 && diffb1.Cmp(diffb2) == 1 && diffb3.Cmp(diffb2) == 1
}
