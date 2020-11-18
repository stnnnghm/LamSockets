package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"net/http"
)

func handleRequest(ctx context.Context, event events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	rc := event.RequestContext
	switch rk := rc.RouteKey; rk {
	case "$connect":
		// manage connect event
		err := connectionStore.AddConnectionID(ctx, rc.ConnectionID)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
			}, err
		}
		break
	case "$disconnect":
		// manage disconnect event
		err := connectionStore.MarkConnectionIDDisconnected(ctx, rc.ConnectionID)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
			}, err
		}
		break
	case "$default":
		// manage every message sent by the clients
		log.Println("Default", rc.ConnectionID)
		err := echo(ctx, event, connectionStore)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode:        http.StatusInternalServerError,
			}, err
		}
		break
	default:
		log.Fatalf("Unknown RouteKey %v", rk)
	}
	// API Gateway is expecting an "everything is okay" answer unless something happens
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}