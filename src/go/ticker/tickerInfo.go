package main

import (
  "context"
  luno "github.com/luno/luno-go"
	"github.com/luno/luno-go/decimal"
  "time"
  "github.com/360EntSecGroup-Skylar/excelize"
  "strconv"
  "fmt"
)

// currencies:
// xth/btc
// xrp/btc
// xrp/zar
// bch/btc
// ltc/btc
// btc/zar

func getTickerRequest(pair string) (*luno.Client, *luno.GetTickerRequest){
  lunoClient := luno.NewClient()
  lunoClient.SetAuth("bgh3q5ss8yvj9", "Q9AqNYS4a4ke-4Yxw_JwtbUbih2mFjITm40luA6Y4M0")

  return lunoClient, &luno.GetTickerRequest{Pair: pair}
}

func getTicker() (decimal.Decimal, decimal.Decimal, luno.Time) {
  res, err := client.GetTicker(context.Background(), reqPointer)
  if err != nil{
    panic(err)
  }
  return res.Ask, res.Bid, res.Timestamp
}

// Global Variables
var client *luno.Client
var reqPointer *luno.GetTickerRequest
var ask decimal.Decimal
var bid decimal.Decimal

func main(){

  fmt.Println("Started")
  
  pairs := []string{"XBTGBP","ETHXBT","XRPXBT","XRPZAR","BCHXBT","LTCXBT","XBTZAR"}
  columns := []string{"A","B","C","D","E","F","G","H","I","J","K","L","M","N"}

  f := excelize.NewFile()
  f.SetCellValue("Sheet1", "A1", "XBTGBP bid")
  f.SetCellValue("Sheet1", "B1", "XBTGBP ask")
  f.SetCellValue("Sheet1", "C1", "ETHXBT bid")
  f.SetCellValue("Sheet1", "D1", "ETHXBT ask")
  f.SetCellValue("Sheet1", "E1", "XRPXBT bid")
  f.SetCellValue("Sheet1", "F1", "XRPXBT ask")
  f.SetCellValue("Sheet1", "G1", "XRPZAR bid")
  f.SetCellValue("Sheet1", "H1", "XRPZAR ask")
  f.SetCellValue("Sheet1", "I1", "BCHXBT bid")
  f.SetCellValue("Sheet1", "J1", "BCHXBT ask")
  f.SetCellValue("Sheet1", "K1", "LTCXBT bid")
  f.SetCellValue("Sheet1", "L1", "LTCXBT ask")
  f.SetCellValue("Sheet1", "M1", "XBTZAR bid")
  f.SetCellValue("Sheet1", "N1", "XBTZAR ask")


  row := 2
  for i := 0; i < 3600; i++{
    // To check progress of ticker
    if i % 60 == 0{
      fmt.Println("Hour: ", i % 60)
    }

    for index,pair := range pairs {

      client, reqPointer = getTickerRequest(pair)
      client.SetTimeout(time.Minute)

      ask, bid, _ = getTicker()
      cell1 := columns[index*2] + strconv.Itoa(row)
      cell2 := columns[index*2+1] + strconv.Itoa(row)

      f.SetCellValue("Sheet1", cell1, bid.String())
      f.SetCellValue("Sheet1", cell2, ask.String())

      time.Sleep(5*time.Second)
    }
    row +=1
    time.Sleep(25*time.Second)
  }



  if err := f.SaveAs("tickerData.xlsx"); err != nil {
    println(err.Error())
  }
  fmt.Println("Ended")

}
