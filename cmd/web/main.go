package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/amtux/user-api-golang/pkg/auth"
	"github.com/amtux/user-api-golang/pkg/models"
	"github.com/go-pg/pg/v10"
	"go.uber.org/zap"
)

type app struct {
	logger *zap.SugaredLogger
	jwt    *auth.JWT
	userDB models.UserDB
}

func main() {
	// parse args
	addr := flag.String("addr", ":8000", "HTTP network address")
	jwtSecret := flag.String("jwt-secret", "samplesecret", "JWT Secret")
	dsn := flag.String("dsn", "postgresql://postgres:postgres@pgdb/postgres?sslmode=disable", "Postgres data source")
	flag.Parse()

	// setup logger
	l, err := zap.NewProduction()
	defer l.Sync()
	logger := l.Sugar()
	if err != nil {
		panic(fmt.Sprintf("failed to start logger %e", err))
	}

	// setup JWT auth
	jwt := auth.JWT{}.New(*jwtSecret)

	// setup DB conn
	opt, err := pg.ParseURL(*dsn)
	if err != nil {
		logger.Fatalw("Failed to parse postgres DSN", err)
	}
	db := pg.Connect(opt)
	ctx := context.Background()
	if err := db.Ping(ctx); err != nil {
		logger.Fatalw("Failed to connect to Database", err)
	}
	defer db.Close()

	// instantiate app
	a := &app{
		logger: logger,
		jwt:    jwt,
		userDB: models.UserDB{DB: db},
	}

	// run migrations
	models.CreateSchema(&a.userDB)

	// start server
	srv := &http.Server{
		Addr:         *addr,
		Handler:      a.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	logger.Infow("Starting server", "URL", addr)
	err = srv.ListenAndServe()
	logger.Fatal(err)
}
