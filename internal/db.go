package internal

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"os"
)

var Db *pgxpool.Pool

func SetupDb() {
	var err error
	Db, err = pgxpool.Connect(context.Background(), os.Getenv("PROTOSVR_DB_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
}
