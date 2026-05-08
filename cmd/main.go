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
	cfg, err := config.NewConfig(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_HOST"), os.Getenv("TOKEN_KEY"))
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

	txManager := repository.NewTxManager(pool)

	userRepository := repository.NewUserRepository(pool)
	userService := service.NewUserService(userRepository, logger)
	userHandler := handler.NewUserHandler(userService, logger)
	userRouter := handler.UserRouter(userHandler)

	currencyRepository := repository.NewCurrencyRepository(pool)
	currencyService := service.NewCurrencyService(currencyRepository)
	currencyHandler := handler.NewCurrencyHandler(currencyService, logger)
	currencyRouter := handler.CurrencyRouter(currencyHandler)

	transactionRepository := repository.NewTransactionRepository(pool)
	transactionService := service.NewTransactionService(transactionRepository)
	transactionHandler := handler.NewTransactionHandler(transactionService, logger)
	transactionRouter := handler.TransactionRouter(transactionHandler)

	accountRepository := repository.NewAccountRepository(pool)
	accountService := service.NewAccountService(accountRepository)
	accountHandler := handler.NewAccountHandler(accountService, logger)
	accountRouter := handler.AccountRouter(accountHandler)

	transferRepository := repository.NewTransferRepository(pool)
	transferService := service.NewTransferService(transferRepository, txManager, accountRepository, transactionRepository)
	transferHandler := handler.NewTransferHandler(transferService, logger)
	transferRouter := handler.TransferRouter(transferHandler)

	withdrawalRepository := repository.NewWithdrawalRepository(pool)
	withdrawalService := service.NewWithdrawalService(withdrawalRepository, txManager, accountRepository, transactionRepository)
	withdrawalHandler := handler.NewWithDrawalHandler(withdrawalService, logger)
	withdrawalRouter := handler.WithdrawalRouter(withdrawalHandler)

	depositRepository := repository.NewDepositRepository(pool)
	depositService := service.NewDepositService(depositRepository, accountRepository, txManager, transactionRepository)
	depositHandler := handler.NewDepositHandler(depositService, logger)
	depositRouter := handler.DepositRouter(depositHandler)

	authRepository := repository.NewAuthRepository(pool)
	authService := service.NewAuthService(&cfg.TokenKey, authRepository, userService, txManager)
	authHandler := handler.NewAuthHandler(userService, authService, logger)
	authRouter := handler.AuthRouter(authHandler)

	root := http.NewServeMux()

	root.Handle("/auth/", http.StripPrefix("/auth", authRouter))

	root.Handle("/users/", http.StripPrefix("/users", userRouter))
	root.Handle("/currencies/", http.StripPrefix("/currencies", currencyRouter))
	root.Handle("/transfers/", http.StripPrefix("/transfers", transferRouter))
	root.Handle("/accounts/", http.StripPrefix("/accounts", accountRouter))
	root.Handle("/withdrawals/", http.StripPrefix("/withdrawals", withdrawalRouter))
	root.Handle("/deposits/", http.StripPrefix("/deposits", depositRouter))
	root.Handle("/transactions/", http.StripPrefix("/transactions", transactionRouter))

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
