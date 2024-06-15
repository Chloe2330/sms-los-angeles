package sms

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

var err error

func SubscriptionWorkflow(ctx workflow.Context, smsDetails SMSDetails) error {
	duration := 12 * time.Second
	logger := workflow.GetLogger(ctx)

	logger.Info("Subscription created", "RecipientPhoneNumber", smsDetails.RecipientPhoneNumber)

	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
		WaitForCancellation: true,
	})

	logger.Info("Sending welcome message...", "RecipientPhoneNumber", smsDetails.RecipientPhoneNumber)
	smsDetails.SubscriptionCount++
	data := SMSDetails{
		TwilioPhoneNumber:    smsDetails.TwilioPhoneNumber,
		RecipientPhoneNumber: smsDetails.RecipientPhoneNumber,
		Message:              "Welcome! You have signed up!",
		IsSubscribed:         true,
		SubscriptionCount:    smsDetails.SubscriptionCount,
	}

	// send welcome message
	err = workflow.ExecuteActivity(ctx, SendMessage, data).Get(ctx, nil)
	if err != nil {
		return err
	}

	// start subscription period
	for smsDetails.IsSubscribed {
		smsDetails.SubscriptionCount++

		logger.Info("Sending subscription message...", "RecipientPhoneNumber", smsDetails.RecipientPhoneNumber)

		data := SMSDetails{
			TwilioPhoneNumber:    smsDetails.TwilioPhoneNumber,
			RecipientPhoneNumber: smsDetails.RecipientPhoneNumber,
			Message:              "This is the recurring message for subscribers",
			IsSubscribed:         true,
			SubscriptionCount:    smsDetails.SubscriptionCount,
		}

		// send subscription messages
		err = workflow.ExecuteActivity(ctx, SendMessage, data).Get(ctx, nil)
		if err != nil {
			return err
		}

		// sleep the workflow until the next subscription message needs to be sent
		if err = workflow.Sleep(ctx, duration); err != nil {
			return err
		}
	}
	return nil
}
