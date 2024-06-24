package sms

import (
	"errors"
	"time"

	"go.temporal.io/sdk/workflow"
)

func SubscriptionWorkflow(ctx workflow.Context, smsDetails SMSDetails) error {
	duration := 30 * time.Second
	logger := workflow.GetLogger(ctx)

	logger.Info("Subscription created", "RecipientPhoneNumber", smsDetails.RecipientPhoneNumber)

	// Query handler
	err := workflow.SetQueryHandler(ctx, "GetDetails", func() (SMSDetails, error) {
		return smsDetails, nil
	})

	if err != nil {
		return err
	}

	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
		WaitForCancellation: true,
	})

	defer func() {

		// create disconnected new context to send cancellation message as async operation
		newCtx, cancel := workflow.NewDisconnectedContext(ctx)

		// clean up resources associated with new context, called when func() exits
		defer cancel()

		// current context (ctx) is canceled
		if errors.Is(ctx.Err(), workflow.ErrCanceled) {
			data := SMSDetails{
				TwilioPhoneNumber:    smsDetails.TwilioPhoneNumber,
				RecipientPhoneNumber: smsDetails.RecipientPhoneNumber,
				Message:              "Your subscription has been canceled. Sorry to see you go!",
				IsSubscribed:         false,
				MessageCount:         smsDetails.MessageCount,
			}
			// send cancellation message
			err := workflow.ExecuteActivity(newCtx, SendMessage, data).Get(newCtx, nil)
			if err != nil {
				logger.Error("Failed to send cancellation message", "Error", err)
			} else {
				// Cancellation received.
				logger.Info("Sent cancellation message", "PhoneNumber", smsDetails.RecipientPhoneNumber)
			}
		}
	}()

	logger.Info("Sending welcome message...", "RecipientPhoneNumber", smsDetails.RecipientPhoneNumber)
	smsDetails.MessageCount++
	data := SMSDetails{
		TwilioPhoneNumber:    smsDetails.TwilioPhoneNumber,
		RecipientPhoneNumber: smsDetails.RecipientPhoneNumber,
		Message:              "Welcome! You have signed up!",
		IsSubscribed:         true,
		MessageCount:         smsDetails.MessageCount,
	}

	// send welcome message
	err = workflow.ExecuteActivity(ctx, SendMessage, data).Get(ctx, nil)
	if err != nil {
		return err
	}

	// start subscription period
	for smsDetails.IsSubscribed {
		smsDetails.MessageCount++

		logger.Info("Sending subscription message...", "RecipientPhoneNumber", smsDetails.RecipientPhoneNumber)

		data := SMSDetails{
			TwilioPhoneNumber:    smsDetails.TwilioPhoneNumber,
			RecipientPhoneNumber: smsDetails.RecipientPhoneNumber,
			Message:              "This is the recurring message for subscribers",
			IsSubscribed:         true,
			MessageCount:         smsDetails.MessageCount,
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
