# SMS Subscription Workflow (Temporal) 

## Set up local development environment for Temporal and Go 

1. Install [Go](https://go.dev/doc/install)

2. Initialize a Go project in a new directory: `go mod init`

3. Install the Temporal Go SDK: `https://go.dev/doc/install`

4. Download the [Temporal CLI](https://go.dev/doc/install)

5. Start a local Temporal service in a new terminal window: `temporal server start-dev`

## Set up Twilio account and secrets 

1. Sign up for [Twilio SMS](https://www.twilio.com/en-us/messaging/channels/sms)

2. Create `.env` file in `server` directory and set Twilio phone number: `TWILIO_PHONE_NUMBER="+1XXXXXXXXX"`

3. Create `.env` file in `worker` directory and set Twilio auth tokens: `TWILIO_ACCOUNT_SID="XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"` `TWILIO_AUTH_TOKEN="XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"`

## Run the workflow 

1. Clone the repository: `git clone https://github.com/Chloe2330/sms-los-angeles.git`

2. Start the server: `cd server` `go run main.go`

3. Start the worker: `cd worker` `go run main.go`

## Send requests to the server

`cd sms-los-angeles`

1. `/subscribe` to daily messages, start workflow execution

    `curl -X POST \  
  http://localhost:4000/subscribe \
  -H 'Content-Type: application/json' \
  -d '{
    "phonenumber": "+1XXXXXXXXX"
  }'`

2. `/unsubscribe` from daily messages, cancel workflow execution

    `curl -X DELETE \                                          
  http://localhost:4000/unsubscribe \
  -H 'Content-Type: application/json' \
  -d '{
    "phonenumber": "+1XXXXXXXXX"
  }'`