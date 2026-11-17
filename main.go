package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	grpcserver "example.com/car-service/grpc"
	pb "example.com/car-service/pb"
	servicehttp "example.com/car-service/http"
	"example.com/car-service/implementation"
	"example.com/car-service/repository"
	"example.com/car-service/transport"
)

func main() {
	httpAddr := flag.String("http.addr", ":8080", "HTTP listen address")
	grpcAddr := flag.String("grpc.addr", ":9090", "gRPC listen address")
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = level.NewFilter(logger, level.AllowDebug())
		logger = log.With(logger,
			"svc", "car",
			"ts", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}
	_ = level.Info(logger).Log("msg", "service started")
	defer func() { _ = level.Info(logger).Log("msg", "service ended") }()

	repo := repository.NewMemoryRepository()
	svc := implementation.NewService(repo, logger)
	endpoints := transport.MakeEndpoints(svc)

	errs := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	// HTTP-транспорт на :8080
	go func() {
		handler := servicehttp.NewService(endpoints, logger)
		_ = logger.Log("transport", "HTTP", "addr", *httpAddr)
		server := &http.Server{
			Addr:    *httpAddr,
			Handler: handler,
		}
		errs <- server.ListenAndServe()
	}()

	// gRPC-транспорт на :9090
	go func() {
		listener, err := net.Listen("tcp", *grpcAddr)
		if err != nil {
			errs <- err
			return
		}
		server := grpc.NewServer()
		pb.RegisterCarServiceServer(server, grpcserver.NewGRPCServer(endpoints, logger))
		reflection.Register(server)
		_ = logger.Log("transport", "gRPC", "addr", *grpcAddr)
		errs <- server.Serve(listener)
	}()

	_ = level.Error(logger).Log("exit", <-errs)
}
