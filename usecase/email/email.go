package email

import (
	"fmt"

	"github.com/go-gomail/gomail"
)

type EmailService struct {
	// Initialize the email service with any required configuration
}

func (e *EmailService) SendInterviewEmail(candidateEmail string, interviewDetails string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "rosequinzley@gmail.com")
	m.SetHeader("To", candidateEmail)
	m.SetHeader("Subject", "Interview Invitation")
	m.SetBody("text/html", fmt.Sprintf("Dear Candidate,<br><br>You have been invited for an interview.<br><br>%s<br><br>Best regards,<br>Your Interview Team", interviewDetails))

	d := gomail.NewDialer("smtp.gmail.com", 587, "rosequinzley@gmail.com", "xhetubsxmjhbdjnr")
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
