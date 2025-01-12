package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"go-todo-app-clean-arch/infrastructure/database"
	"go-todo-app-clean-arch/infrastructure/web"
	"go-todo-app-clean-arch/pkg"
	"go-todo-app-clean-arch/pkg/logger"
)

func main() {
	appEnv := pkg.GetEnvDefault("APP_ENV", "development")
	if appEnv == "development" {
		err := godotenv.Load(".env.development")
		if err != nil {
			logger.Warn("Error loading .env.local file")
		}
		logger.Info("Loaded SECRET from .env.development: ", os.Getenv("SECRET"))
	}

	db, err := database.NewDatabaseSQLFactory(database.InstanceMySQL)
	if err != nil {
		logger.Fatal(err.Error())
	}

	// server, err := web.NewServer(web.InstanceGin, db)
	server, err := web.NewServer(web.InstanceEcho, db)
	if err != nil {
		logger.Fatal(err.Error())
	}
	go func() {
		if err := server.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal(err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")
	defer logger.Sync()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error(fmt.Sprintf("Server Shutdown: %s", err.Error()))
	}
	<-ctx.Done()
}
