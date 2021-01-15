package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"time"

	"github.com/Stardome-Team/Service-Template/internal/account"
	"github.com/Stardome-Team/Service-Template/internal/config"
	"github.com/Stardome-Team/Service-Template/internal/errors"
	"github.com/Stardome-Team/Service-Template/pkg/database"
	"github.com/Stardome-Team/Service-Template/pkg/logset"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator"
	_ "github.com/lib/pq"
)

var version = "0.0.1"
var flagConfig = flag.String("config", "./config/local.yml", "Path to configuration file")

const (
	// passwordRequirementRegexPattern requires the password to have at least on lowercase letter, uppercase letter and one number
	passwordRequirementRegexPattern = `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d]{8,}$`
)

func main() {
	flag.Parse()

	// create application logger tag with application version
	logger := logset.New().With(nil, "version", version)

	// load application configurations
	cfg, err := config.Load(*flagConfig)

	if err != nil {
		logger.Errorf("Failed to load application configuration: %s", err)
		os.Exit(1)
	}

	// connect to the database
	db, err := database.OpenConnection(cfg)

	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	srv := &http.Server{
		Addr:    ":1010",
		Handler: buildHandler(cfg, db, logger),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logger.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorf("Server Shutdown: %s", err)
	}
	logger.Info("Server exiting")
}

func buildHandler(cfg *config.Config, db *database.DB, logger logset.Logger) *gin.Engine {
	router := gin.Default()

	router.Use(errors.ErrorHandlerMiddleware())

	group := router.Group("/api")

	account.CreateHandlers(group, account.NewService(account.NewRepository(db, logger), logger), logger)

	return router
}

func registerValidations() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("validatepassword", passwordRegexValidator)
	}
}

// passwordRegexValidator validator function for password regex
var passwordRegexValidator validator.Func = func(fl validator.FieldLevel) bool {
	password, ok := fl.Field().Interface().(string)

	if ok {
		if matched, err := regexp.MatchString(passwordRequirementRegexPattern, password); err == nil {
			return matched
		}
	}

	return false
}
