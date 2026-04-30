package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alanzhumalin/bank/internal/config"
	"github.com/alanzhumalin/bank/internal/db"
	"github.com/alanzhumalin/bank/internal/handler"
	"github.com/alanzhumalin/bank/internal/logger"
	"github.com/alanzhumalin/bank/internal/repository"
	"github.com/alanzhumalin/bank/internal/service"
	"github.com/joho/godotenv"
)

func main() {

	logger := logger.NewLogger()

	//connect to db, we need to get config from environment
	godotenv.Load()
	cfg, err := config.NewConfig(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_HOST"))

	if err != nil {
		logger.Fatal().Err(err).Msg("Error occured while getting environment files")
		os.Exit(1)
	}

	//postgresql://username:password@localhost:5432/dbname?sslmode=disable

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:5432/%s?sslmode=disable", cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.DbName)

	pool, err := db.ConnectDB(dsn)

	if err != nil {
		logger.Fatal().Err(err).Msg("Error occured while connecting to db")
	}

	userRepository := repository.NewUserRepository(pool)
	userService := service.NewUserService(userRepository, logger)
	userHandler := handler.NewUserHandler(userService, logger)
	userRouter := handler.UserRouter(userHandler)

	currencyRepository := repository.NewCurrencyRepository(pool)
	currencyService := service.NewCurrencyService(currencyRepository)
	currencyHandler := handler.NewCurrencyHandler(currencyService, logger)
	currencyRouter := handler.CurrencyRouter(currencyHandler)

	transferRepository := repository.NewTransferRepository(pool)
	transferService := service.NewTransferService(transferRepository)
	transferHandler := handler.NewTransferHandler(transferService, logger)
	transferRouter := handler.TransferRouter(transferHandler)

	accountRepository := repository.NewAccountRepository(pool)
	accountService := service.NewAccountService(accountRepository)
	accountHandler := handler.NewAccountHandler(accountService, logger)
	accountRouter := handler.AccountRouter(accountHandler)

	root := http.NewServeMux()

	root.Handle("/users/", http.StripPrefix("/users", userRouter))
	root.Handle("/currency/", http.StripPrefix("/currency", currencyRouter))
	root.Handle("/transfer/", http.StripPrefix("/transfer", transferRouter))
	root.Handle("/account/", http.StripPrefix("/account", accountRouter))

	srv := http.Server{
		Addr:    ":8081",
		Handler: root,
	}

	go func() {
		logger.Info().Str("Port", srv.Addr).Msg("Server started")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal().Err(err).Msg("Error occured while starting the server")
			os.Exit(1)
		}
	}()

	ch := make(chan os.Signal, 1)

	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	<-ch

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error().Err(err).Msg("Error occured")
	}

}
