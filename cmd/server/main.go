package main

import (
	"context"
	"log/slog"

	"github.com/freekobie/hazel/handlers"
	"github.com/freekobie/hazel/mail"
	"github.com/freekobie/hazel/postgres"
	"github.com/freekobie/hazel/services"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	cfg := loadConfig()

	db, err := pgxpool.New(context.Background(), cfg.PostgresURL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping(context.Background())
	if err != nil {
		panic(err)
	}

	mailer := mail.NewMailer(cfg.MailConfig)
	handler := handlers.NewHandler(services.NewUserService(postgres.NewUserStore(db), mailer))

	app := &application{
		handler: handler,
	}

	slog.Info("Starting server")
	err = app.start()
	if err != nil {
		panic(err)
	}
}
