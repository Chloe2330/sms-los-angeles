package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sms"

	"github.com/joho/godotenv"
	"go.temporal.io/sdk/client"
)

var temporalClient client.Client

type RequestData struct {
	PhoneNumber string `json:"phonenumber"`
}

type ResponseData struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func main() {
	var err error

	// create client
	temporalClient, err = client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})

	if err != nil {
		panic(err)
	}

	fmt.Printf("Starting the web server on %s\n", sms.ClientHostPort)

	http.HandleFunc("/subscribe", subscribeHandler)
	http.HandleFunc("/unsubscribe", unsubscribeHandler)
	_ = http.ListenAndServe(":4000", nil)
}

// create subscribe handler, which collects the subscriber's phone number
func subscribeHandler(w http.ResponseWriter, r *http.Request) {

	// only respond to POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// ensure JSON request
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Invalid Content-Type, expecting application/json", http.StatusUnsupportedMediaType)
		return
	}

	var requestData RequestData

	// decode request into variable
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Error processing request body", http.StatusBadRequest)
		return
	}

	// check if the phone number is blank
	if requestData.PhoneNumber == "" {
		http.Error(w, "Phone number is blank", http.StatusBadRequest)
		return
	}

	// use the phone number as the id in the workflow.
	workflowOptions := client.StartWorkflowOptions{
		ID:                                       requestData.PhoneNumber,
		TaskQueue:                                sms.TaskQueueName,
		WorkflowExecutionErrorWhenAlreadyStarted: true,
	}

	// Define the SMSDetails struct
	subscription := sms.SMSDetails{
		TwilioPhoneNumber:    GetEnvVar("TWILIO_PHONE_NUMBER"),
		RecipientPhoneNumber: requestData.PhoneNumber,
		Message:              "Welcome to the Subscription Workflow!",
		IsSubscribed:         true,
		SubscriptionCount:    0,
	}

	// Execute the Temporal Workflow to start the subscription.
	_, err = temporalClient.ExecuteWorkflow(context.Background(), workflowOptions, sms.SubscriptionWorkflow, subscription)

	if err != nil {
		http.Error(w, "Couldn't sign up user. Please try again.", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	// build response
	responseData := ResponseData{
		Status:  "success",
		Message: "Signed up.",
	}

	// send headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 Created status code

	// send response
	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		log.Print("Could not encode response JSON", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// create unsubscribe handler, accessed at localhost:4000/unsubscribe
func unsubscribeHandler(w http.ResponseWriter, r *http.Request) {

	// only respond to DELETE
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// ensure JSON request
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Invalid Content-Type, expecting application/json", http.StatusUnsupportedMediaType)
		return
	}

	var requestData RequestData

	// decode request into variable
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Error processing request body", http.StatusBadRequest)
		return
	}

	// check if the phone number is blank
	if requestData.PhoneNumber == "" {
		http.Error(w, "Phone number is blank", http.StatusBadRequest)
		return
	}
	workflowID := requestData.PhoneNumber

	// cancel and return a CancelledError to the Workflow Execution
	err = temporalClient.CancelWorkflow(context.Background(), workflowID, "")
	if err != nil {
		http.Error(w, "Couldn't unsubscribe. Please try again.", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	// build response
	responseData := ResponseData{
		Status:  "success",
		Message: "Unsubscribed.",
	}

	// send headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted) // 202 Accepted status code

	// send response
	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		log.Print("Could not encode response JSON", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// use godot package to load/read the .env file and return Twilio secrets
func GetEnvVar(key string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
