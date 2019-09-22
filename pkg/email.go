package pkg

import gomail "gopkg.in/gomail.v2"

//WahaEmail 发送电子邮件
type WahaEmail struct {
	UserName string
	PassWord string
	Host     string
	Port     int
	IsHTML   bool
}

//NewEmail 新建email对象
func NewEmail(name string, pwd string, host string, port int, isHTML bool) *WahaEmail {
	return &WahaEmail{
		UserName: name,
		PassWord: pwd,
		Host:     host,
		Port:     port,
		IsHTML:   isHTML,
	}
}

//Send 发送电子邮件
func (email *WahaEmail) Send(to string, title string, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", email.UserName)
	m.SetHeader("To", to)
	//m.SetAddressHeader("Cc", "test@example.com", "Dan")
	m.SetHeader("Subject", title)
	var contentType string

	if email.IsHTML {
		contentType = "text/html"
	} else {
		contentType = "text/plain"
	}
	m.SetBody(contentType, body)
	d := gomail.NewDialer(email.Host, email.Port, email.UserName, email.PassWord)
	err := d.DialAndSend(m)
	return err
}
