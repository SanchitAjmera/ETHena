package main

import (
    "fmt"
    "net/smtp"
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
    // Sender data.
    from := "profit.profit.profit.icl@gmail.com"
    password := "Password123??"
    // Receiver email address.
    to := []string{
        "profit.profit.profit.icl@gmail.com",
    }
    // smtp server configuration.
    smtpServer := smtpServer{host: "smtp.gmail.com", port: "587"}
    // Message.
    var messageStr string
    // add customised switch cases here
    switch action{
      case "BUY":
        fallthrough
      case "SELL":
        messageStr = "Congratulation! Sanchit's bot made a " + action + " at " + info1.String()
      case "START":
        messageStr = "NEWS! Sanchit's bot has begun trading"
    }
    message := []byte(messageStr)
    // Authentication.
    auth := smtp.PlainAuth("", from, password, smtpServer.host)
    // Sending email.
    err := smtp.SendMail(smtpServer.Address(), auth, from, to, message)

    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println("Email Sent!")
}
