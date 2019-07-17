package emailer

import (
	"fmt"
	"github.com/scorredoira/email"
	"net/mail"
	"net/smtp"
)

func EmailFile(filepath, targetAddress, fromEmail, fromPassword string) error {
	// compose the message
	m := email.NewMessage("Hi", "this is the body")
	m.From = mail.Address{Name: "Scanner", Address: fromEmail}
	m.To = []string{targetAddress}

	// add attachments
	if err := m.Attach(filepath); err != nil {
		return fmt.Errorf("could not attach file: %s", err)
	}

	// send it
	auth := smtp.PlainAuth("", fromEmail, fromPassword, "smtp.gmail.com")
	if err := email.Send("smtp.gmail.com:587", auth, m); err != nil {
		return fmt.Errorf("could not send email: %s", err)
	}
	return nil
}

//func authEmail() *gmail.Service {
//	data, err := ioutil.ReadFile(keyFilePath)
//	if err != nil {
//		log.Fatal("could not read key file to authenticate with email server", err)
//	}
//	conf, err := google.JWTConfigFromJSON(data, "https://www.googleapis.com/auth/gmail.compose")
//	if err != nil {
//		log.Fatal("could not auth with the gmail api", err)
//	}
//	client := conf.Client(context.TODO())
//	gmailService, err := gmail.New(client)
//	if err != nil {
//		log.Fatal("could not establish new gmail service based on the authenticated client", err)
//	}
//	return gmailService
//}
