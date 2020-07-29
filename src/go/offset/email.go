package main

import (
    "fmt"
    "gopkg.in/gomail.v2" // perform go get <this> when initialising servers
    "github.com/luno/luno-go/decimal"
)

type smtpServer struct {
 host string
 port string
}
// Address URI to smtp server
func (s *smtpServer) Address() string {
 return s.host + ":" + s.port
}

func email(action string, info1 decimal.Decimal, info2 decimal.Decimal) {
  // Sender/Reciever data.
  m := gomail.NewMessage()
  m.SetHeader("From", "profit.profit.profit.icl@gmail.com")
  m.SetHeader("To", "profit.profit.profit.icl@gmail.com") // can add multiple recievers
  var messageStr string
  // add customised switch cases here
  switch action{
    case "BUY":
      messageStr = "NEWS! Sanchit's bot placed a" + action + " order at " + info1.String()
    case "SELL":
      messageStr = "NEWS! Sanchit's bot placed a" + action + " order at " + info1.String()
    case "START":
      messageStr = "NEWS! Sanchit's bot has begun trading"
    case "BOUGHT":
      messageStr = "NEWS! Sanchit's bot completed the BUY order at " + info1.String()
    case "SOLD":
      messageStr = "NEWS! Sanchit's bot completed the SELL order at " + info1.String()
  }

  m.SetHeader("Subject", messageStr)
  m.SetBody("text/html", "")

  d := gomail.NewDialer("smtp.gmail.com", 587, "profit.profit.profit.icl@gmail.com", "Password123??")

	if err := d.DialAndSend(m); err != nil {
    fmt.Println("ERROR! Update email NOT successfully sent")
    return
	}

  fmt.Println("Update email successfully sent")
}
