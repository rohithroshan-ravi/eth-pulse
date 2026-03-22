package config

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func ConnectDB() {
	connString := "postgres://postgres:Chottu374@8220@localhost:5432/eth_pulse"
	var err error
	DB, err = pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatalf("Mamea! DB kulla poga mudiyala: %v\n", err)
	}

	err = DB.Ping(context.Background())
	if err != nil {
		log.Fatal("Database Connection Failed, Check Password!")
	}

	fmt.Println("✅ Database Connected Successfully (Pool Mode)!")
}
