package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/alanzhumalin/bank/internal/cache"
	"github.com/alanzhumalin/bank/internal/config"
	"github.com/alanzhumalin/bank/internal/db"
	"github.com/alanzhumalin/bank/internal/handler"
	"github.com/alanzhumalin/bank/internal/logger"
	"github.com/alanzhumalin/bank/internal/middleware"
	"github.com/alanzhumalin/bank/internal/repository"
	"github.com/alanzhumalin/bank/internal/service"
	"github.com/joho/godotenv"
)

func main() {

	logger := logger.NewLogger()

	//connect to db, we need to get config from environment
	godotenv.Load()

	cfg, err := config.NewConfig(os.Getenv("REDIS_PASSWORD"), os.Getenv("REDIS_ADDRESS"), os.Getenv("REDIS_DB"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_HOST"), os.Getenv("TOKEN_KEY"))
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

	defer pool.Close()

	redisDbInt, _ := strconv.Atoi(cfg.RedisDb)

	redisClient, err := db.InitRedisClient(cfg.RedisAddress, cfg.RedisPassword, redisDbInt)

	if err != nil {
		logger.Fatal().Err(err).Msg("Error occured while initializing redis client")
		os.Exit(1)
	}

	defer redisClient.Close()

	duration := 1 * time.Minute
	limitReq := int64(10)

	rateLimiter, err := cache.NewRateLimiter(redisClient, limitReq, duration)
	rateLimitMiddlewareFactory, err := middleware.NewRateLimiterMiddleware(rateLimiter)

	if err != nil {
		logger.Fatal().Err(err).Msg("Error occured with initializing the middleware")
	}
	rateLimitMiddleware := rateLimitMiddlewareFactory.RateLimiterMiddleware()

	rbac := middleware.NewRbacMiddleware().RBAC

	txManager := repository.NewTxManager(pool)

	idempotencyRepo := repository.NewIdempotencyRepo(pool)

	idempotencyRedis, err := cache.NewIdempotencyStore(redisClient, time.Duration(1*time.Minute), time.Duration(30*time.Minute), time.Duration(1*24*time.Hour))

	if err != nil {
		logger.Fatal().Err(err).Msg("Error occured initializing the idempotency redis")
		os.Exit(1)
	}
	blackListTokenRedis, err := cache.NewTokenBlackList(redisClient)

	if err != nil {
		logger.Fatal().Err(err).Msg("Error occured initializing the black token redis")
		os.Exit(1)
	}
	authMiddleware := middleware.NewAuthMiddleWare(&cfg.TokenKey, blackListTokenRedis).Middleware()

	userRepository := repository.NewUserRepository(pool)
	userService := service.NewUserService(userRepository, logger)
	userHandler := handler.NewUserHandler(userService, logger)
	userRouter := handler.UserRouter(userHandler, authMiddleware, rbac)

	currencyRepository := repository.NewCurrencyRepository(pool)
	currencyService := service.NewCurrencyService(currencyRepository)
	currencyHandler := handler.NewCurrencyHandler(currencyService, logger)
	currencyRouter := handler.CurrencyRouter(currencyHandler, authMiddleware, rbac)

	transactionRepository := repository.NewTransactionRepository(pool)
	transactionService := service.NewTransactionService(transactionRepository)
	transactionHandler := handler.NewTransactionHandler(transactionService, logger)
	transactionRouter := handler.TransactionRouter(transactionHandler, authMiddleware, rateLimitMiddleware, rbac)

	accountRepository := repository.NewAccountRepository(pool)
	accountService := service.NewAccountService(accountRepository)
	accountHandler := handler.NewAccountHandler(accountService, logger)
	accountRouter := handler.AccountRouter(accountHandler, authMiddleware, rbac)

	transferRepository := repository.NewTransferRepository(pool)
	transferService := service.NewTransferService(transferRepository, txManager, accountRepository, transactionRepository)
	transferHandler := handler.NewTransferHandler(transferService, logger)
	transferRouter := handler.TransferRouter(transferHandler, authMiddleware)

	withdrawalRepository := repository.NewWithdrawalRepository(pool)
	withdrawalService := service.NewWithdrawalService(idempotencyRepo, withdrawalRepository, txManager, accountRepository, transactionRepository)
	withdrawalHandler := handler.NewWithDrawalHandler(idempotencyRedis, withdrawalService, logger)
	withdrawalRouter := handler.WithdrawalRouter(withdrawalHandler, authMiddleware)

	depositRepository := repository.NewDepositRepository(pool)
	depositService := service.NewDepositService(depositRepository, accountRepository, txManager, transactionRepository)
	depositHandler := handler.NewDepositHandler(depositService, logger)
	depositRouter := handler.DepositRouter(depositHandler)

	authRepository := repository.NewAuthRepository(pool)
	authService := service.NewAuthService(blackListTokenRedis, &cfg.TokenKey, authRepository, userService, txManager)
	authHandler := handler.NewAuthHandler(userService, authService, logger)
	authRouter := handler.AuthRouter(authHandler, authMiddleware)

	root := http.NewServeMux()

	root.Handle("/auth/", http.StripPrefix("/auth", authRouter))

	root.Handle("/deposits/", http.StripPrefix("/deposits", depositRouter))

	root.Handle("/users/", http.StripPrefix("/users", userRouter))
	root.Handle("/currencies/", http.StripPrefix("/currencies", currencyRouter))
	root.Handle("/transfers/", http.StripPrefix("/transfers", transferRouter))
	root.Handle("/accounts/", http.StripPrefix("/accounts", accountRouter))
	root.Handle("/withdrawals/", http.StripPrefix("/withdrawals", withdrawalRouter))
	root.Handle("/transactions/", http.StripPrefix("/transactions", transactionRouter))

	srv := http.Server{
		Addr:    ":8081",
		Handler: withCORS(root),
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

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
