package main

import (
	"context"
	"flag"
	"github.com/Scarlet-Fairy/gateway/gateway"
	"log"
)

var (
	managerEndpoint    = flag.String("manager-endpoint", "localhost:8081", "manager gRPC server endpoint")
	sloweaterEndpoint  = flag.String("sloweater-endpoint", "localhost:8082", "sloweater gRPC server endpoint")
	logWatcherEndpoint = flag.String("log-watcher-endpoint", "localhost:8085", "log-watcher gRPC server endpoint")
	gatewayUrl         = flag.String("http-addr", ":8080", "addr to bind the server")
)

func main() {
	flag.Parse()

	ctx := context.Background()
	opts := gateway.Options{
		Address: *gatewayUrl,
		Endpoints: struct {
			Manager    gateway.Endpoint
			LogWatcher gateway.Endpoint
		}{
			Manager: gateway.Endpoint{
				Address: *managerEndpoint,
			},
			LogWatcher: gateway.Endpoint{
				Address: *logWatcherEndpoint,
			},
		},
		OpenApiDir: "gen/openAPI",
	}

	if err := gateway.Run(ctx, opts); err != nil {
		log.Fatal(err)
	}
}
