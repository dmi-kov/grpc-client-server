package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/grpc-client-server/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"io"
	"net/http"
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
	stream, err := c.CallURL(context.Background(), &api.URLMessage{Url: *strURL})
	if err != nil {
		logger.Fatalf("Error when calling server: %s", err)
	}

	bodyAcc := make([]byte, 0)
	headers := make([]byte, 0)

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Fatalf("Fail reading stream: %s", err)
			break
		}
		bodyAcc = append(bodyAcc, resp.Response...)
		headers = append(headers, resp.Headers...)
	}

	var headersMap http.Header
	err = json.Unmarshal(headers, &headersMap)
	if err != nil {
		logger.Fatalf("Fail Unmarshal headers: %s", err)
	}

	logger.Infof("Response from server: headers %v \n body \n %v", headersMap, string(bodyAcc))
}
