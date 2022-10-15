package myUtils

import (
	"bytes"
	"crypto/tls"
	"gopkg.in/gomail.v2"
	"html/template"
	"talkRoom/setting"
)

var (
	dialer    *gomail.Dialer
	emailConf = setting.Config.Email
	templ     *template.Template
)

func init() {
	dialer = gomail.NewDialer(emailConf.Host, emailConf.Port, emailConf.UserName, emailConf.Token)
	// 关闭SSL协议认证
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	var err error
	if templ, err = template.ParseFiles("templates/temp/verificationCode.html"); err != nil {
		panic(err)
	}
}

func SentVerificationCode(email string, vCode string) error {
	var buff bytes.Buffer
	if err := templ.Execute(&buff, vCode); err != nil {
		return err
	}
	m := gomail.NewMessage()
	m.SetHeader("From", emailConf.UserName)                // 发件人
	m.SetHeader("To", email)                               // 收件人，可以多个收件人，但必须使用相同的 SMTP 连接
	m.SetHeader("Subject", "Register & verification code") // 邮件主题
	m.SetBody("text/html", buff.String())

	return dialer.DialAndSend(m)
}
