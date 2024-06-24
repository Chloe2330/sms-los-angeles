package sms

const TaskQueueName string = "sms_subscription"
const ClientHostPort string = "localhost:4000"

type SMSDetails struct {
	TwilioPhoneNumber    string `json:"twilioPhoneNumber"`
	RecipientPhoneNumber string `json:"recipientPhoneNumber"`
	Message              string `json:"message"`
	IsSubscribed         bool   `json:"isSubscribed"`
	MessageCount         int    `json:"messageCount"`
}
