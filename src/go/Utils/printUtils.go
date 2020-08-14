package Utils

import(
  "fmt"
  "time"
  "os"
  "os/exec"
  "github.com/luno/luno-go/decimal"
)

var StartDay time.Time

// function to print the information about the bot

func clear(){
   cmd := exec.Command("clear")
   cmd.Stdout = os.Stdout
   cmd.Run()
}

func LoadScreen(){
  for i:= 0; i < 52; i++ {
    clear()
    fmt.Printf("           ")
    fmt.Printf("         ")
    printEthena()
    fmt.Printf("\n\n")
    fmt.Printf("              ")
    for j := 0; j < i; j ++ {
      fmt.Printf("#")
    }
    for j:= 50; j >= i; j -- {
      fmt.Printf("-")
    }
    time.Sleep(time.Second /15)
  }
  fmt.Printf("\n")
}

func PrintStatus(b *RsiBot, currAsk decimal.Decimal, currBid decimal.Decimal, status string, values []([]decimal.Decimal)){
  clear()
  fmt.Printf("Date | %v", StartDay.Format("02 Jan 06"))
  fmt.Printf("\n")
  fmt.Printf("Time | %v", StartDay.Format("15:04:05 MST"))
  fmt.Printf("\n")
  fmt.Printf("Pair | ETHBTX                     ")
  fmt.Printf("Status:")
  fmt.Printf("\n")
  fmt.Printf("User | %v            ", User)
  for i := 0; i < (20-len(status)); i++{
    fmt.Printf(" ")
  }
  fmt.Printf("%v",status)
  fmt.Printf("\n")
  fmt.Printf("Bid  | %v", currBid)
  fmt.Printf("\n")
  fmt.Printf("Ask  | %v", currAsk)
  if values != nil {
    printIndicators(values[0], values[1], values[2], values[3])
  }

  if b != nil {
    fmt.Printf("\n\n")
    fmt.Printf("                 ")
    fmt.Printf("Buy Price  | %v", b.BuyPrice)
    fmt.Printf("\n")
    fmt.Printf("                 ")
    fmt.Printf("Sell Price | %v", b.SellPrice)
  }

}

func printEthena(){
  fmt.Printf("\n\n\n\n\n\n\n")
  fmt.Println("                    #####  #####  #   #  #####  #   #    #")
  fmt.Println("                    #        #    #   #  #      ##  #   # #")
  fmt.Println("                    #####    #    #####  #####  # # #  #####")
  fmt.Println("                    #        #    #   #  #      #  ##  #   #")
  fmt.Println("                    #####    #    #   #  #####  #   #  #   #")

}
func printIndicators(rsiValues []decimal.Decimal,macdValues []decimal.Decimal, stopLossValues []decimal.Decimal, scoreValues []decimal.Decimal) {
  one := decimal.NewFromInt64(1)
  fmt.Printf("\n\n\n\n")
  fmt.Printf("                 ")
  fmt.Printf("|    RSI   |   MACD   | STOPLOSS |  SCORE   |\n")
  fmt.Printf("                 ")
  fmt.Printf("|          |          |          |          |\n")
  for i := 0 ; i < len(rsiValues) - 1 ; i++{
    fmt.Printf("                 ")
    rsiValue := rsiValues[i].Div(one, 16).String()
    macdValue := macdValues[i].Div(one, 16).String()
    scoreValue := scoreValues[i].Div(one, 16).String()
    stopLossValue := stopLossValues[i].Div(one, 16).String()
    fmt.Printf("| %v",rsiValue[:8])
    fmt.Printf(" | ")
    fmt.Printf("%v",macdValue[:8])
    fmt.Printf(" | ")
    fmt.Printf("%v",stopLossValue[:8])
    fmt.Printf(" | ")
    fmt.Printf("%v",scoreValue[:8])
    fmt.Printf(" | ")
    fmt.Printf("\n")
  }
  fmt.Printf("               > ")
  rsiValue := rsiValues[len(rsiValues) - 1].Div(one, 16).String()
  macdValue := macdValues[len(rsiValues) - 1].Div(one, 16).String()
  scoreValue := scoreValues[len(rsiValues) - 1].Div(one, 16).String()
  stopLossValue := stopLossValues[len(rsiValues) - 1].Div(one, 16).String()
  fmt.Printf("| %v",rsiValue[:8])
  fmt.Printf(" | ")
  fmt.Printf("%v",macdValue[:8])
  fmt.Printf(" | ")
  fmt.Printf("%v",stopLossValue[:8])
  fmt.Printf(" | ")
  fmt.Printf("%v",scoreValue[:8])
  fmt.Printf(" | ")
  fmt.Printf("\n")

}

// func main(){
//   loadScreen()
// }
