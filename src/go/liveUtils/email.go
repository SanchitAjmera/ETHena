package liveUtils

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

func Email(action string, yield decimal.Decimal) {
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
      case "GRAPH":
        messageStr = "Daily update: "
        if yield.Sign() == 1 {
          messageStr += "PROFIT! £"
        } else {
          messageStr += "LOSS! £"
        }
        messageStr += yield.String()
      case "START":
        messageStr = "NEWS! Your bot has begun trading"
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
    fmt.Println("Update email sent successfully")
}
