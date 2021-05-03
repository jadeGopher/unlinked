package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"net/http"
	"os"
	"unlinked/cmd/handlers"
	"unlinked/proto"
)

func main() {
	var (
		logger *zap.Logger
		err    error
	)
	if logger, err = zap.NewProduction(); err != nil {
		panic(err)
	}

	var db *sql.DB
	if db, err = sql.Open("postgres", os.Getenv("database")); err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}

	var mux = runtime.NewServeMux()
	if err = proto.RegisterUnlinkedServiceHandlerServer(
		context.Background(),
		mux,
		handlers.New(db, logger),
	); err != nil {
		panic(err)
	}

	_ = http.ListenAndServe(fmt.Sprintf(":%s", "3000"), mux)
}
