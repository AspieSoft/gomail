package gomail

import (
	"bufio"
	"os"
	"strings"
	"testing"
)

var emailAuth map[string]string = map[string]string{}

func getEmailAuth(){
	file, err := os.Open("./test.auth")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		val := strings.SplitN(scanner.Text(), ":", 2)
		if len(val) != 2 {
			continue
		}
		emailAuth[strings.TrimSpace(val[0])] = strings.TrimSpace(val[1])
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func TestMail(t *testing.T) {
	getEmailAuth()

	mailer, err := NewMailer(emailAuth["email"], emailAuth["pwd"], HOST.Gmail)
	if err != nil {
		t.Error("mailer", err)
	}

	err = mailer.Send([]string{emailAuth["to"]}, "Test", "<h3>This Is A Test</h3>")
	if err != nil {
		t.Error("send", err)
	}

	err = mailer.SendFrom(emailAuth["from"], []string{emailAuth["to"]}, "Test", "<h3>This Is A Test</h3>")
	if err != nil {
		t.Error("send", err)
	}
}
