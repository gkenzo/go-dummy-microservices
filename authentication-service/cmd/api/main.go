package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gkenzo/go-dummy-microservices/authentication-service/cmd/data"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type Config struct {
	DB     *sql.DB
	Models data.Models
}

const (
	webPort = "80"
)

func main() {
	conn := connectToDB()
	if conn == nil {
		log.Panic("Can't connect to database!")
	}

	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Starting service on port %s\n", webPort)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dbConnsRetries := 0
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready...")
			dbConnsRetries++
		} else {
			log.Println("Successfully connected to the database.")
			return connection
		}

		if dbConnsRetries > 10 {
			log.Println(err)
			return nil
		}

		time.Sleep(2 * time.Second)
	}

}
