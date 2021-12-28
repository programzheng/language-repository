package email

import (
	"bytes"
	"html/template"
	"log"
	"net/smtp"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

//Request struct
type Request struct {
	from    string
	to      []string
	subject string
	body    string
}

func newRequest(from string, to []string, subject, body string) *Request {
	return &Request{
		from:    from,
		to:      to,
		subject: subject,
		body:    body,
	}
}

func (r *Request) sendEmail() (bool, error) {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + r.subject + "!\n"
	from := "From: " + r.from + "\n"
	to := "To: " + strings.Join(r.to, ",") + "\n"
	msg := []byte(subject + from + to + mime + "\n" + r.body)
	addr := os.Getenv("MAIL_HOST") + ":" + os.Getenv("MAIL_PORT")

	auth := smtp.PlainAuth("", os.Getenv("MAIL_USERNAME"), os.Getenv("MAIL_PASSWORD"), os.Getenv("MAIL_HOST"))

	if err := smtp.SendMail(addr, auth, r.from, r.to, msg); err != nil {
		return false, err
	}
	return true, nil
}

func getTemplatePath(fileName string) string {
	path := filepath.Join(basepath, "../dist/email/"+fileName)
	return path
}

func (r *Request) parseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}

func SendEmailByHtml(to []string, emailFileName string, title string, subject string, data map[string]interface{}) bool {
	r := newRequest(os.Getenv("MAIL_FROM"), to, title, subject)

	err := r.parseTemplate(getTemplatePath(emailFileName), data)
	if err != nil {
		log.Fatal(err)
	}

	ok, err := r.sendEmail()
	if err != nil {
		log.Fatal(err)
	}

	return ok
}
