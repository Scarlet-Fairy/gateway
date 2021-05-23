package gateway

import (
	"context"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	gw "github.com/Scarlet-Fairy/gateway/pb" // Update
)

type Endpoint struct {
	Address string
}

type Options struct {
	Address   string
	Endpoints struct {
		Manager    Endpoint
		LogWatcher Endpoint
	}
	OpenApiDir string
	Mux        []runtime.ServeMuxOption
}

func Run(ctx context.Context, options Options) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	gwMux := runtime.NewServeMux(options.Mux...)
	opts := []grpc.DialOption{grpc.WithInsecure()}

	if err := gw.RegisterManagerHandlerFromEndpoint(
		ctx,
		gwMux,
		options.Endpoints.Manager.Address,
		opts,
	); err != nil {
		return err
	}

	if err := gw.RegisterLogWatcherHandlerFromEndpoint(
		ctx,
		gwMux,
		options.Endpoints.LogWatcher.Address,
		opts,
	); err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/swagger/", openAPIServer(options.OpenApiDir))
	mux.Handle("/api/", gwMux)

	s := http.Server{
		Addr:    options.Address,
		Handler: allowCORS(mux),
	}

	go func() {
		<-ctx.Done()
		log.Println("Shutting down the http server")
		if err := s.Shutdown(context.Background()); err != nil {
			log.Fatalf("Failed to shutdown http server: %v", err)
		}
	}()

	log.Printf("Starting listening at %s\n", options.Address)
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}
