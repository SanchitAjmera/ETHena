package liveUtils

import (
    "fmt"
    "github.com/luno/luno-go/decimal"
    "gopkg.in/gomail.v2" // perform go get <this> when initialising servers
)


func Email(action string, yield decimal.Decimal) {
  m := gomail.NewMessage()
	m.SetHeader("From", "profit.profit.profit.icl@gmail.com")
	m.SetHeader("To", "profit.profit.profit.icl@gmail.com") // can add multiple recievers
  var messageStr string

  switch action{
    case "GRAPH":
      messageStr = "Daily update: "
      if yield.Sign() == 1 {
        messageStr += "PROFIT! £" + yield.String()
      } else if yield.Sign() == 1{
        yield.Mul(decimal.NewFromFloat64(-1,8))
        messageStr += "LOSS! £" + yield.String()
      } else {
        messageStr += "FLAT! £" + yield.String()
      }
      m.Attach("../main/graph.png")
      messageStr += yield.String()
    case "START":
      messageStr = "NEWS! Your bot has begun trading"
  }

	m.SetHeader("Subject", messageStr)
	m.SetBody("text/html", ".")

	d := gomail.NewDialer("smtp.gmail.com", 587, "profit.profit.profit.icl@gmail.com", "Password123??")

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

  fmt.Println("Update email successfully sent!")
}
