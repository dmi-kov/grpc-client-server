package main

import (
	"context"
	"github.com/grpc-client-server/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	listenPort = ":8081"
	once       = &sync.Once{}
	ctx        context.Context
	cancelFunc context.CancelFunc
)

func main() {
	zapLogger, _ := zap.NewDevelopment()
	defer zapLogger.Sync()
	logger := zapLogger.Sugar()

	ctx, _ := getApplicationContext()

	if err := listenAndServe(ctx, listenPort, logger); err != nil {
		logger.Fatalf("failed to start server: %v", err)
	}
}

func listenAndServe(ctx context.Context, port string, logger *zap.SugaredLogger) (err error) {
	grpcServer := grpc.NewServer()

	go func() {
		logger.Infof("server listening on port %v", port)
		lis, err := net.Listen("tcp", port)
		if err != nil {
			logger.Errorf("failed to listen: %v", err)
		}
		if err := grpcServer.Serve(lis); err != nil {
			logger.Errorf("failed to serve GRPC: %v", err)
		}
	}()

	s := api.Handler{}
	api.RegisterAPIServer(grpcServer, &s)

	<-ctx.Done()

	logger.Info("stopping server")
	grpcServer.GracefulStop()

	return
}

func getApplicationContext() (context.Context, context.CancelFunc) {
	once.Do(func() {
		ctx, cancelFunc = context.WithCancel(context.Background())

		go func() {
			signals := []os.Signal{syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT}
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, signals...)
			defer signal.Reset(signals...)
			<-sigChan
			cancelFunc()
		}()
	})

	return ctx, cancelFunc
}
