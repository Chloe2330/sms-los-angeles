package sms

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sms/metro"
	"strconv"
	"strings"

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

func GetMetroInfo(ctx context.Context) (string, error) {

	swiftlyKey := GetEnvVar("SWIFTLY_API_KEY")
	latitude := GetEnvVar("LATITUDE")
	longitude := GetEnvVar("LONGITUDE")

	url := fmt.Sprintf("https://api.goswift.ly/real-time/lametro-rail/predictions-near-location?lat=%s&lon=%s", latitude, longitude)

	// Create a new HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	// Set the headers
	req.Header.Set("Accept", "application/json, application/xml")
	req.Header.Set("Authorization", swiftlyKey)

	// Make the GET request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	// Check the HTTP status code
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Received non-200 response code: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	var metroPredictions metro.MetroPredictions
	if err := json.Unmarshal(body, &metroPredictions); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}

	var dataSlice []metro.MetroPredictionsFormatted

	// extract useful information from json response
	for _, predictionsData := range metroPredictions.Data.PredictionsData {
		for _, dest := range predictionsData.Destinations {
			var minsSlice []string
			for _, prediction := range dest.Predictions {
				minsSlice = append(minsSlice, strconv.Itoa(prediction.Min))
			}
			mins := strings.Join(minsSlice, ", ")

			data := metro.MetroPredictionsFormatted{
				RouteName:        predictionsData.RouteName,
				StopName:         predictionsData.StopName,
				DestStopName:     dest.Headsign,
				MinsUntilArrival: mins,
			}
			dataSlice = append(dataSlice, data)
		}
	}

	// concatenate strings with string builder for efficiency
	var builder strings.Builder
	for _, data := range dataSlice {
		fmt.Fprintf(&builder, "\nRoute: %s\nClosest Stop: %s\nDest: %s\nMinutes Until Arrival: %s\n\n", data.RouteName, data.StopName, data.DestStopName, data.MinsUntilArrival)
	}
	text := builder.String()

	return text, err
}

// use godot package to load/read the .env file and return Twilio secrets
func GetEnvVar(key string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
