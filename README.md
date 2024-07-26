# SMS Subscription Workflow (Temporal) 

## Set up local development environment for Temporal and Go 

1. Install [Go](https://go.dev/doc/install)

2. Initialize a Go project in a new directory: `go mod init`

3. Install the Temporal Go SDK: `https://github.com/temporalio/sdk-go`

4. Download the [Temporal CLI](https://learn.temporal.io/getting_started/go/dev_environment/#set-up-a-local-temporal-service-for-development-with-temporal-cli) 

5. Start a local Temporal service in a new terminal window: `temporal server start-dev`

## Set up Twilio account and secrets 

1. Sign up for [Twilio SMS](https://www.twilio.com/en-us/messaging/channels/sms)

2. Create `.env` file in `server` directory and set Twilio phone number: `TWILIO_PHONE_NUMBER="+1XXXXXXXXXX"`

3. Create `.env` file in `worker` directory and set Twilio auth tokens: `TWILIO_ACCOUNT_SID="XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"` `TWILIO_AUTH_TOKEN="XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"`

## Run the workflow 

1. Clone the repository: `git clone https://github.com/Chloe2330/sms-los-angeles.git`

2. Start the server: `cd server` `go run main.go`

3. Start the worker: `cd worker` `go run main.go`

## Send requests to the server

`cd sms-los-angeles`
### Endpoints 

1. `/subscribe` to daily messages, start workflow execution
```bash
curl -X POST \
  http://localhost:4000/subscribe \
  -H 'Content-Type: application/json' \
  -d '{
    "phonenumber": "+1XXXXXXXXXX"
  }'
```

2. `/unsubscribe` from daily messages, cancel workflow execution
```bash
 curl -X DELETE 'http://localhost:4000/unsubscribe?phonenumber=1XXXXXXXXXX'
```

3. `/details` to query information about an existing workflow
```bash
 curl 'http://localhost:4000/details?phonenumber=1XXXXXXXXXX'
```