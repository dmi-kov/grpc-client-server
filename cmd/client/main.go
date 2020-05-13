package main

import (
	"context"
	"flag"
	"github.com/grpc-client-server/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net/url"
)

var (
	strURL     = flag.String("url", "", "valid url")
	targetPort = ":8081"
)

func main() {
	zapLogger, _ := zap.NewDevelopment()
	defer zapLogger.Sync()
	logger := zapLogger.Sugar()
	flag.Parse()

	_, err := url.ParseRequestURI(*strURL)
	if err != nil {
		logger.Fatalf("passed URL not valid: %v", err)
	}

	conn, err := grpc.Dial(targetPort, grpc.WithInsecure())
	if err != nil {
		logger.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := api.NewAPIClient(conn)
	resp, err := c.CallURL(context.Background(), &api.URLMessage{Url: *strURL})
	if err != nil {
		logger.Fatalf("Error when calling server: %s", err)
	}

	received, err := resp.Recv()
	if err != nil {
		logger.Fatalf("Fail receive message from GRPC stream: %s", err)
	}

	var first100 string
	bodyString := string(received.Body)
	if len(bodyString) > 100 {
		first100 = bodyString[0:100]
	} else {
		first100 = bodyString
	}
	logger.Infof("Response from server: \n <BODY LEN>: %v \n <BODY>: %s \n <HEADERS>: %v", len(first100), first100, string(received.Headers))
}
