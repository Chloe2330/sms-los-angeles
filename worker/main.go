package main

import (
	"log"
	"sms"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	// create client and worker
	c, err := client.Dial(client.Options{
		HostPort:  client.DefaultHostPort,
		Namespace: client.DefaultNamespace,
	})
	if err != nil {
		log.Fatalln("Unable to create Temporal Client.", err)
	}
	defer c.Close()

	// create Worker
	w := worker.New(c, sms.TaskQueueName, worker.Options{})

	// register Activity and Workflow
	w.RegisterWorkflow(sms.SubscriptionWorkflow)
	w.RegisterActivity(sms.SendMessage)
	w.RegisterActivity(sms.GetMetroInfo)
	w.RegisterActivity(sms.GetCoordinates)

	log.Println("Worker is starting.")

	// Listen to Task Queue
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start Worker.", err)
	}
}
