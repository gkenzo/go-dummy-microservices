package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gkenzo/go-dummy-microservices/log-service/cmd/api/data"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	defaultWebPort    = "80"
	rpcPort           = "5001"
	gRpcPort          = "50001"
	defaultDbUser     = "admin"
	defaultDbPassword = "password"
	defaultDbURL      = "bd-logger-mongo://mongo:27017"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	// connect to logger db (mongo)

	dbClient, err := connectToMongo()
	if err != nil {
		log.Panicf("an error occurred while connecting to database %v", err)
	}
	client = dbClient

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	app.serve()
}

func (app *Config) serve() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultDbUser
	}
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic()
	}
}

func connectToMongo() (*mongo.Client, error) {
	url, user, password := getDbCredentials()
	connOpts := options.Client().ApplyURI(url)
	connOpts.SetAuth(options.Credential{
		Username: user,
		Password: password,
	})
	c, err := mongo.Connect(context.TODO(), connOpts)

	return c, err
}

// Gets db credentials from environment. if none provided, use default values.
func getDbCredentials() (url string, username string, password string) {
	url, uSet := os.LookupEnv("DB_URL")
	if !uSet {
		url = defaultDbURL
	}
	u, uSet := os.LookupEnv("DB_USER")
	if !uSet {
		u = defaultDbUser
	}
	p, pSet := os.LookupEnv("DB_PASSWORD")
	if !pSet {
		p = defaultDbPassword
	}

	return url, u, p
}
