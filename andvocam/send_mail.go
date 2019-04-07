package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
	"strconv"
	"strings"

	"github.com/tkanos/gonfig"
)

type configuration struct {
	Address    string
	Password   string
	SMTPServer string
	SMTPPort   int
	Subject    string
}

func sendEmail(clientEmail string, comment string) {

	// file, err := os.Open("config.cfg")
	configuration := configuration{}
	confError := gonfig.GetConf("config.cfg", &configuration)
	if confError != nil {
		fmt.Println(confError)
	}
	// Set up authentication information.
	auth := smtp.PlainAuth("",
		configuration.Address,
		configuration.Password,
		configuration.SMTPServer)

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	from := mail.Address{Name: "Онлайн конвертер", Address: configuration.Address}
	to := mail.Address{Name: "Онлайн конвертер", Address: configuration.Address}
	title := configuration.Subject
	body := clientEmail + " оставил комментарий: \r\n" + comment + ".\r\n"

	header := make(map[string]string)
	header["From"] = from.Address
	header["To"] = to.Address
	header["Subject"] = encodeRFC2047(title)
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	err := smtp.SendMail(configuration.SMTPServer+":"+strconv.Itoa(configuration.SMTPPort), auth, configuration.Address, []string{to.Address}, []byte(message))
	if err != nil {
		log.Fatal(err)
	}
}
func encodeRFC2047(String string) string {
	// use mail's rfc2047 to encode any string
	addr := mail.Address{String, ""}
	return strings.Trim(addr.String(), " <>")
}
