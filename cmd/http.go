package cmd

import (
	"fmt"
	"github.com/saleh-ghazimoradi/X-Gopher/config"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/gateway/middlewares"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/gateway/routes"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/server"
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

		middleware := middlewares.NewMiddleware(cfg, logger)

		register := routes.NewRegister(
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
