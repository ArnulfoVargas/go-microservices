package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

var counts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting Auth service")
	// TODO connect to db
	conn := ConnectToDB()
	if conn == nil {
		panic("Cant connect to postgres")
	}

	// Config
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.Routes(),
	}

	err := srv.ListenAndServe()

	if err != nil {
		panic(err.Error())
	}
}

func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		println(err.Error())
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		println(err.Error())
		return nil, err
	}

	return db, nil
}

func ConnectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := OpenDB(dsn)

		if err != nil {
			log.Println("Postgres not yet ready")
			counts++
		} else {
			log.Println("Connected to postgres")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for two seconds")
		time.Sleep(2 * time.Second)
		continue
	}
}
