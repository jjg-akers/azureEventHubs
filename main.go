package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	aad "github.com/Azure/azure-amqp-common-go/aad"
	eventhubs "github.com/Azure/azure-event-hubs-go"
)

func main() {

	//create an authorization provider for your Event Hubs client that uses these credentials:
	tokenProvider, err := aad.NewJWTProvider(aad.JWTProviderWithEnvironmentVars())
	if err != nil {
		log.Fatalf("failed to configure AAD JWT provider: %s\n", err)
	}

	//Create Event Hubs client
	hub, err := eventhubs.NewHub("namespaceName", "hubName", tokenProvider)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer hub.Close(ctx)
	if err != nil {
		log.Fatalf("failed to get hub %s\n", err)
	}

	//code to send messages
	//In the following snippet, use (1) to send messages interactively from a terminal,
	// or (2) to send messages within your program:

	// 1. send messages at the terminal
	ctx = context.Background()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("Input a message to send: ")
		text, _ := reader.ReadString('\n')
		hub.Send(ctx, eventhubs.NewEventFromString(text))
	}

	// 2. send messages within program
	ctx = context.Background()
	hub.Send(ctx, eventhubs.NewEventFromString("hello Azure!"))

	//Get the IDs of the partitions in your event hub:
	info, err := hub.GetRuntimeInformation(ctx)
	if err != nil {
		log.Fatalf("failed to get runtime info: %s\n", err)
	}
	log.Printf("got partition IDs: %s\n", info.PartitionIDs)

}
