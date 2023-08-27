package gomail

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"net/smtp"

	"github.com/thanhpk/randstr"
)

type mime struct {
	Text string
	Html string
	Png string
	Jpeg string
	Jpg string
	Mp4 string
	Mp3 string
	Json string
	Msword string
}

var MIME mime = mime{
	Text: "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";",
	Html: "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";",
	Png: "MIME-version: 1.0;\nContent-Type: image/png; charset=\"UTF-8\";",
	Jpeg: "MIME-version: 1.0;\nContent-Type: image/jpeg; charset=\"UTF-8\";",
	Jpg: "MIME-version: 1.0;\nContent-Type: image/jpg; charset=\"UTF-8\";",
	Mp4: "MIME-version: 1.0;\nContent-Type: video/mp4; charset=\"UTF-8\";",
	Mp3: "MIME-version: 1.0;\nContent-Type: audio/mp3; charset=\"UTF-8\";",
	Json: "MIME-version: 1.0;\nContent-Type: application/json; charset=\"UTF-8\";",
	Msword: "MIME-version: 1.0;\nContent-Type: application/msword; charset=\"UTF-8\";",
}

type MailHost struct {
	host string
	port string
}

type Host struct {
	Gmail MailHost
}

var HOST Host = Host{
	Gmail: MailHost{host: "smtp.gmail.com", port: "587"},
}

type Mailer struct {
	Send func(to []string, subject string, body string, mime ...string) error
	SendFrom func(from string, to []string, subject string, body string, mime ...string) error
}

var bytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}
var encKey string

func init(){
	encKey = randstr.String(32)
}

func NewHost(host string, port string) MailHost {
	return MailHost{host, port}
}

func NewMailer(email string, pwd string, host MailHost, fromEmail ...string) (Mailer, error) {
	encPWD, err := encrypt(pwd, encKey)
	if err != nil {
		return Mailer{}, err
	}

	auth := map[string]string{
		"email": email,
		"pwd": encPWD,
		"host": host.host,
		"port": host.port,
	}

	fromDef := email
	if len(fromEmail) != 0 {
		fromDef = fromEmail[0]
	}

	send := func(to []string, subject string, body string, mime ...string) error {
		mimeType := MIME.Html
		if len(mime) != 0 {
			mimeType = mime[0]
		}
		return sendEmail(auth, fromDef, to, subject, body, mimeType)
	}

	sendFrom := func(from string, to []string, subject string, body string, mime ...string) error {
		mimeType := MIME.Html
		if len(mime) != 0 {
			mimeType = mime[0]
		}
		return sendEmail(auth, from, to, subject, body, mimeType)
	}

	return Mailer{send, sendFrom}, nil
}

func sendEmail(authData map[string]string, fromName string, to []string, subject string, body string, mime string) error {
	// sender data
	pwd, err := decrypt(authData["pwd"], encKey)
	if err != nil {
		return err
	}
	from := authData["email"] // example@gmail.com
	password := pwd // abcdefghijklmnop

	// smtp
	host := authData["host"] // smtp.gmail.com
	port := authData["port"] // 587
	address := host + ":" + port

	// message
	toList := ""
	for i, email := range to {
		toList += email
		if i != len(to) - 1 {
			toList += ", "
		}
	}
	msg := []byte("From: " + fromName + "\nTo: " + toList + "\nSubject: " + subject + "\n" + mime + "\n\n" + body)

	// auth
	auth := smtp.PlainAuth("", from, password, host)

	// send email
	err = smtp.SendMail(address, auth, from, to, msg)
	if err != nil {
		return err
	}
	return nil
}

func encrypt(text string, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, bytes)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func decrypt(text string, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	cipherText, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", err
	}
	cfb := cipher.NewCFBDecrypter(block, bytes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}
