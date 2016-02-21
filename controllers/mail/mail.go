package mail

import (
	"fmt"
	"net/smtp"
	"path/filepath"
	"strings"

	"github.com/astaxie/beego/config"
)

// Mailer define smtp server and email account info
type Mailer struct {
	account  string
	password string
	host     string
	port     string
}

// New create Mailer
// account email account name
// pwd password
// host smtp server host
// port smtp server port default: 25
func New(account, pwd, host, port string) (mail *Mailer) {
	return &Mailer{account: account, password: pwd, host: host, port: port}
}

// NewServiceMailer create Mailer host service@coderpage.com
func NewServiceMailer() (mail *Mailer, err error) {

	confPath := filepath.Join("conf", "mail.conf")
	conf, err := config.NewConfig("ini", confPath)
	if err != nil {
		return nil, err
	}
	account := conf.String("serviceAccount")
	pass := conf.String("servicePass")
	host := conf.String("serviceHost")
	port := conf.String("servicePort")
	if port == "" {
		port = "25"
	}
	mailer := &Mailer{account: account, password: pass, host: host, port: port}
	return mailer, nil
}

// SendMail send email
// Header e.g
// To: 820631390@qq.com
// From: uauth <service@coderpage.com>
// Subject: subject
// Content-Type: text/html; charset=UTF-8
// body
func (mailer *Mailer) SendMail(to, name, subject, mailtype, body string) (err error) {
	auth := smtp.PlainAuth("", mailer.account, mailer.password, mailer.host)
	var contentType string
	if mailtype == "html" {
		contentType = "Content-Type: text/html; charset=UTF-8"
	} else {
		contentType = "Content-Type: text/plain; charset=UTF-8"
	}

	msg := fmt.Sprintf("From: %s<%s>\r\nTo: %s\r\nSubject: %s\r\n%s\r\n\r\n%s", name, mailer.account, to, subject, contentType, body)
	sendTo := strings.Split(to, ";")
	addr := mailer.host + ":" + mailer.port
	err = smtp.SendMail(addr, auth, mailer.account, sendTo, []byte(msg))
	return err
}
