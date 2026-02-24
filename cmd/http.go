package cmd

import (
	"context"
	"fmt"
	"github.com/saleh-ghazimoradi/X-Gopher/config"
	"github.com/saleh-ghazimoradi/X-Gopher/infra/mongodb"
	"github.com/saleh-ghazimoradi/X-Gopher/infra/redis"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/gateway/handlers"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/gateway/middlewares"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/gateway/routes"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/repository"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/server"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/service"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

// httpCmd represents the http command
var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("http called")

		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

		cfg, err := config.GetInstance()
		if err != nil {
			logger.Error("Failed to load configuration", "error", err)
			os.Exit(1)
		}

		redisDB := redis.NewRedis(
			redis.WithHost(cfg.Redis.Host),
			redis.WithPort(cfg.Redis.Port),
			redis.WithPassword(cfg.Redis.Password),
			redis.WithDB(cfg.Redis.DB),
			redis.WithDialTimeout(cfg.Redis.DialTimeout),
			redis.WithReadTimeout(cfg.Redis.ReadTimeout),
			redis.WithWriteTimeout(cfg.Redis.WriteTimeout),
			redis.WithPoolSize(cfg.Redis.PoolSize),
			redis.WithPoolTimeout(cfg.Redis.PoolTimeout),
		)

		if err := redisDB.Ping(context.Background()); err != nil {
			logger.Error("Failed to connect to redis", "error", err)
		}

		defer func() {
			if err := redisDB.Close(context.Background()); err != nil {
				logger.Error("Failed to close redis", "error", err)
			}
		}()

		//redisClient := redisDB.Connect(context.Background())

		mongo := mongodb.NewMongoDB(
			mongodb.WithHost(cfg.MongoDB.Host),
			mongodb.WithPort(cfg.MongoDB.Port),
			mongodb.WithUser(cfg.MongoDB.User),
			mongodb.WithPass(cfg.MongoDB.Pass),
			mongodb.WithDBName(cfg.MongoDB.DBName),
			mongodb.WithAuthSource(cfg.MongoDB.AuthSource),
			mongodb.WithMaxPoolSize(cfg.MongoDB.MaxPoolSize),
			mongodb.WithMinPoolSize(cfg.MongoDB.MinPoolSize),
			mongodb.WithTimeout(cfg.MongoDB.Timeout),
		)

		client, mongodb, err := mongo.Connect()
		if err != nil {
			logger.Error("Failed to connect to MongoDB", "error", err)
			os.Exit(1)
		}

		defer func() {
			if err := client.Disconnect(context.Background()); err != nil {
				logger.Error("Failed to disconnect from MongoDB", "error", err)
				os.Exit(1)
			}
		}()

		middleware := middlewares.NewMiddleware(cfg, logger)

		tokenRepository := repository.NewTokenRepository(mongodb, "token")
		userRepository := repository.NewUserRepository(mongodb, "user")
		postRepository := repository.NewPostRepository(mongodb, "post")
		commentRepository := repository.NewCommentRepository(mongodb, "comment")
		messageRepository := repository.NewMessageRepository(mongodb, "message")
		unreadMessageRepository := repository.NewUnreadMessageRepository(mongodb, "unreadMessage")
		notificationRepository := repository.NewNotificationRepository(mongodb, "notification")

		authService := service.NewAuthService(cfg, userRepository, tokenRepository)
		userService := service.NewUserService(userRepository, notificationRepository)
		postService := service.NewPostService(userRepository, commentRepository, postRepository, notificationRepository)
		messageService := service.NewMessageService(messageRepository, unreadMessageRepository)
		notificationService := service.NewNotificationService(notificationRepository)

		authHandler := handlers.NewAuthHandler(authService)
		userHandler := handlers.NewUserHandler(userService)
		postHandler := handlers.NewPostHandler(postService)
		messageHandler := handlers.NewMessageHandler(messageService)
		notificationHandler := handlers.NewNotificationHandler(notificationService)

		authRoute := routes.NewAuthRoute(authHandler)
		userRoute := routes.NewUserRoute(middleware, userHandler)
		postRoute := routes.NewPostRoute(middleware, postHandler)
		messageRoute := routes.NewMessageRoute(messageHandler)
		notificationRoute := routes.NewNotificationRoute(notificationHandler)

		register := routes.NewRegister(
			routes.WithAuthRoute(authRoute),
			routes.WithUserRoute(userRoute),
			routes.WithPostRoute(postRoute),
			routes.WithMessageRoute(messageRoute),
			routes.WithNotificationRoute(notificationRoute),
			routes.WithMiddlewares(middleware),
		)

		httpServer := server.NewHTTPServer(
			server.WithHost(cfg.HTTPServer.Host),
			server.WithPort(cfg.HTTPServer.Port),
			server.WithHandler(register.RegisterRoutes()),
			server.WithReadTimeout(cfg.HTTPServer.ReadTimeout),
			server.WithIdleTimeout(cfg.HTTPServer.IdleTimeout),
			server.WithWriteTimeout(cfg.HTTPServer.WriteTimeout),
			server.WithErrorLog(slog.NewLogLogger(logger.Handler(), slog.LevelError)),
			server.WithLogger(logger),
		)

		logger.Info("starting server", "addr", cfg.HTTPServer.Host+":"+cfg.HTTPServer.Port, "env", cfg.Application.Environment)
		if err := httpServer.Connect(); err != nil {
			logger.Error("Failed to connect to http server", "error", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(httpCmd)
}
