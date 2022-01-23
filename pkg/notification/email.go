package notification

import (
	"log"
	"net/smtp"
)

func sendEmail() {
	from := "inwtax@gmail.com"
	password := "0554861712"
	to := []string{"boonyarit.b@securitypitch.com"}
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	message := []byte("Hello Test.")

	// Create authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Send actual message
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		log.Fatal(err)
	}
}
