package sms

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"go.temporal.io/sdk/activity"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"

	"github.com/joho/godotenv"
)

func SendMessage(ctx context.Context, smsInfo SMSDetails) error {
	accountSid := GetEnvVar("TWILIO_ACCOUNT_SID")
	authToken := GetEnvVar("TWILIO_AUTH_TOKEN")

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(smsInfo.RecipientPhoneNumber)
	params.SetFrom(smsInfo.TwilioPhoneNumber)
	params.SetBody(smsInfo.Message)

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		activity.GetLogger(ctx).Info("Unable to send message to subscriber", "RecipientPhoneNumber", smsInfo.RecipientPhoneNumber)
	} else {
		body, _ := json.Marshal(*resp.Body)
		activity.GetLogger(ctx).Info("Successfully sent this message: "+string(body), "RecipientPhoneNumber", smsInfo.RecipientPhoneNumber)
	}
	return err
}

// use godot package to load/read the .env file and return Twilio secrets
func GetEnvVar(key string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
