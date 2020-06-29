package main

import (
        luno "github.com/luno/luno-go"
        "fmt"
)

func buy(pf portfolio, stock decimal.Decimal, price decimal.Decimal) {
  currFunds := pf.funds
  if (currFunds < stock.Mul(price)) {
    fmt.PrintLn("Cannot afford to buy ", stock, " stock at ", price, " price")
  } else {
    pf.funds = pf.funds.Sub(stock.Mul(price))
    pf.stock = pf.stock.Add(stock)
    fmt.PrintLn("Bought ",stock, " stock at ", price, " price")
  }
}

func sell(pf portfolio, stock decimal.Decimal, price decimal.Decimal) {
  currStock := pf.stock
  if (currStock < stock) {
    fmt.PrintLn("Not enough stock to sell ", stock, " stock at ", price, " price")
  } else {
    pf.funds = pf.funds.Add(stock.Mul(price))
    pf.stock = pf.stock.Sub(stock)
    fmt.PrintLn("Sold ",stock, " stock at ", price, " price")
  }
}

type portfolio struct { //Features common to every bot
  funds decimal.Decimal
  stock decimal.Decimal
  tradingPeriod int64 //Trading period in minutes
  currRow int64
}

type SMA_bot struct {
  pf portfolio
  offset decimal.Decimal //Offset size
  numOfDecisions int64 //Length of short moving average as multiple of period
}

func (b SMA_bot) trade() {
  pastBids []decimal.Zero()

  i := 1
  currBid := getBid(b.pf.currRow)

  for i <= b.numOfDecisions * b.pf.tradingPeriod {
    pastBids[i] = getBid(b.pf.currRow - i)
    i++
  }

  mean:= SMA(pastBids)

  if currBid.Cmp(mean.Add(b.offset)) == 1 {
    sell(b.pf, b.pf.stock, currBid)
  } else if currBid.Cmp(mean.Sub(b.offset)) == -1 {
    buy(b.pf, b.pf.stock, getAsk(b.pf.currRow))
  }

  b.pf.currRow += b.numOfDecisions * b.pf.tradingPeriod
}
