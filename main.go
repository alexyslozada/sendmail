package main

import (
	"fmt"
	"log"
	"bytes"
	"net/mail"
	"net/smtp"
	"crypto/tls"
	"html/template"
)

type Dest struct {
	Name string
}

func checkErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func main()  {
	from := mail.Address{"Alexys Envio", "alexysest@gmail.com"}
	to := mail.Address{"Alexys Lozada", "alexyslozada@msn.com"}
	subject := "Enviando correo desde GO"
	dest := Dest{Name: to.Address}

	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subject
	headers["Content-Type"] = `text/html; charset="UTF-8"`

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	t, err := template.ParseFiles("template.html")
	checkErr(err)

	buf := new(bytes.Buffer)
	err = t.Execute(buf, dest)
	checkErr(err)

	message += buf.String()

	servername := "smtp.gmail.com:465"
	host := "smtp.gmail.com"

	auth := smtp.PlainAuth("", "alexysest@gmail.com", "Clave123+", host)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName: host,
	}

	conn, err := tls.Dial("tcp", servername, tlsConfig)
	checkErr(err)

	client, err := smtp.NewClient(conn, host)
	checkErr(err)

	err = client.Auth(auth)
	checkErr(err)

	err = client.Mail(from.Address)
	checkErr(err)

	err = client.Rcpt(to.Address)
	checkErr(err)

	w, err := client.Data()
	checkErr(err)

	_, err = w.Write([]byte(message))
	checkErr(err)

	err = w.Close()
	checkErr(err)

	client.Quit()
}