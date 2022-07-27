# GoMail

[![donation link](https://img.shields.io/badge/buy%20me%20a%20coffee-square-blue)](https://buymeacoffee.aspiesoft.com)

Simplifies Sending Emails In Go.

## Installation

```shell script

  go get github.com/AspieSoft/go-regex

```

## Usage

```go

import (
  "github.com/AspieSoft/gomail/"
)

// creating an authenticated mailer
var mailer gomail.Mailer = gomail.NewMailer(
  "example@gmail.com", // a real email address
  "abcdefghijklmnop", // email password or an app password
  gomail.HOST.gmail, // an email host
  "MyName <noreply@example.com>", // (optional) Custom Name to send emails as by default
  // Note: A custom name Must be a valid alias in gmail or may be required with your host of choice
)

// custom host
var gmailHost gomail.Host = gomail.NewHost("smtp.gmail.com", "587")


// sending an email
func main(){
  mailer.send(
    []string{"user1@example.com", "user2@example.com"}, // list of emails to send to
    "My Email Subject",
    "My Email Body",
    gomail.MIME.html, // (optional) default: html
  )

  mailer.sendFrom(
    "Support <support@example.com>", // change the alias an email is sent from in place of the default
    []string{"user1@example.com", "user2@example.com"},
    "My Email Subject",
    "My Email Body",
    gomail.MIME.text,
  )

  mailer.send(
    []string{"me@example.com"},
    "Test Email",
    "<h1>Hello, Email!</h1>",
    "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";",
  )
}

```
