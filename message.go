package main

import (
	"os"
	"regexp"
	"strings"

	"github.com/go-mail/mail"
)

var rxEmail = regexp.MustCompile(".+@.+\\..+")
var rxPhone = regexp.MustCompile(`((?:\D|^)\d{10})`)
var rxBusyKid = regexp.MustCompile(`^1\d\d\d\d\d\d\d`)

// type DropdownItem struct {
// 	Name  string
// 	Value string
// }

// var fruits = map[string]interface{}{
// 	"Annette": "Annette",
// 	"Cheluna": "Cheluna",
// }

type Message struct {
	Name       string
	DateTime   string
	Phone      string
	Address    string
	BusyKidNum string
	Email      string
	Restaurant string
	Content    string
	Errors     map[string]string
}

func (msg *Message) Validate() bool {
	msg.Errors = make(map[string]string)

	match := rxEmail.Match([]byte(msg.Email))
	if match == false {
		msg.Errors["Email"] = "Please enter a valid email address"
	}

	if strings.TrimSpace(msg.Name) == "" {
		msg.Errors["Name"] = "Please enter your name"
	}

	if strings.TrimSpace(msg.Address) == "" {
		msg.Errors["Address"] = "Please enter your Address"
	}

	if strings.TrimSpace(msg.BusyKidNum) == "" {
		msg.Errors["BusyKidNum"] = "Please enter the Busy Kid Payment Confirmation Number"
	}

	matchBusyKid := rxBusyKid.Match([]byte(msg.BusyKidNum))
	if matchBusyKid == false {
		msg.Errors["BusyKidNum"] = "Please enter a valid BusyKid Confirmation Number"
	}

	if strings.TrimSpace(msg.DateTime) == "" {
		msg.Errors["DateTime"] = "Please enter the date and time of pickup"
	}

	if strings.TrimSpace(msg.Restaurant) == "" {
		msg.Errors["Restaurant"] = "Please enter the Restaurant"
	}

	if strings.TrimSpace(msg.Content) == "" {
		msg.Errors["Content"] = "Please enter details about your order"
	}

	matchPhone := rxPhone.Match([]byte(msg.Phone))
	if strings.TrimSpace(msg.Phone) == "" {
		msg.Errors["Phone"] = "Please enter your phone number"
	}
	if matchPhone == false {
		msg.Errors["Phone"] = "Please enter a valid phone number"
	}

	return len(msg.Errors) == 0
}

func (msg *Message) Deliver() error {
	this_msg := "\n" +
		"Name: " + msg.Name + "\n" +
		"DateTime: " + msg.DateTime + "\n" +
		"Phone: " + msg.Phone + "\n" +
		"Email: " + msg.Email + "\n" +
		"Restaurant: " + msg.Restaurant + "\n" +
		"Content: " + msg.Content

	email := mail.NewMessage()
	email.SetHeader("To", os.Getenv("boazform_emailto"))
	email.SetHeader("From", os.Getenv("boazform_from"))
	email.SetHeader("Reply-To", msg.Email)
	email.SetHeader("Subject", "New message via Contact Form")
	email.SetBody("text/plain", this_msg)

	username := os.Getenv("boazform_username")
	password := os.Getenv("boazform_password")

	return mail.NewDialer("smtp.gmail.com", 465, username, password).DialAndSend(email)
}

func (msg *Message) Deliver_Receipt() error {
	this_msg := "Thank you for your delivery order, we will contact you shortly " +
		"with an update.  The following are the details of your delivery " +
		"order.\n"
	this_msg = this_msg + "\n" +
		"Name: " + msg.Name + "\n" +
		"DateTime: " + msg.DateTime + "\n" +
		"Phone: " + msg.Phone + "\n" +
		"Email: " + msg.Email + "\n" +
		"Restaurant: " + msg.Restaurant + "\n" +
		"Content: " + msg.Content

	email := mail.NewMessage()
	email.SetHeader("To", msg.Email)
	email.SetHeader("From", os.Getenv("boazform_from"))
	//email.SetHeader("Reply-To", msg.Email)
	email.SetHeader("Subject", "Thank you for your business!")
	email.SetBody("text/plain", this_msg)

	username := os.Getenv("boazform_username")
	password := os.Getenv("boazform_password")

	return mail.NewDialer("smtp.gmail.com", 465, username, password).DialAndSend(email)
}
