package server

import (
	"context"
	"errors"
	"funovation_23/graph/generated"
	"funovation_23/internal/config"

	"funovation_23/internal/setup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.uber.org/zap"
)

func Run() error {
	log.Println("Reding configuration...")
	configuration := config.LoadConfig()
	dbConn, err := setup.SetupDb(configuration)
	if err != nil {
		log.Println("Error while connecting to database")
		return err
	}
	httpClient := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
	s3Client, err := setup.SetupS3Client(configuration.S3Config, &httpClient)
	if err != nil {
		return err
	}
	log.Printf("%+v\n", configuration)
	resolver, err := setup.NewResolver(dbConn, *configuration, s3Client)
	if err != nil {
		return err
	}
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: resolver,
	}))

	router := mux.NewRouter()
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Println("Calling serve")
	return serve(router, configuration)
}

func serve(mux *mux.Router, config *config.Config) error {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := corsMiddleware.Handler(mux)
	api := http.Server{
		Addr:         "0.0.0.0:" + config.Server.Port,
		ReadTimeout:  config.Server.ReadTimeout,
		WriteTimeout: config.Server.WriteTimeout,
		Handler:      handler,
	}
	serverErrors := make(chan error, 1)
	go func() {
		sugar.Infof("Connect to http://localhost:%s/ for GraphQL playground", config.Server.Port)
		if config.Server.TLSEnable {
			serverErrors <- api.ListenAndServeTLS(config.Server.TLSCertPath, config.Server.TLSKeyPath)
		} else {
			serverErrors <- api.ListenAndServe()
		}
	}()

	select {
	case err := <-serverErrors:
		return err

	case sig := <-shutdown:
		ctx, cancel := context.WithTimeout(context.Background(), config.Server.ShutdownTimeout)
		defer cancel()

		err := api.Shutdown(ctx)
		if err != nil {
			err = api.Close()
		}

		switch {
		case sig == syscall.SIGKILL:
			return errors.New("integrity error shuting down")

		case err != nil:
			return err
		}
		return nil
	}
}
