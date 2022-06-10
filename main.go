package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"go-minimal-api/auth"
	"go-minimal-api/uploads"
	"go.uber.org/zap"
	"net/http"
	"os"
)

func main() {
	// init logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	sugar.Infow("started my image api")

	lg := auth.Login{Logger: logger}
	err := lg.Setup()
	if err != nil {
		fmt.Println("Error setting up the db:", err.Error())
	}

	store, _ := uploads.NewStorage(context.Background())
	if err != nil {
		fmt.Println("Error setting up the storage:", err.Error())
	}

	e := echo.New()
	e.GET("/ping", pingHandler)
	e.POST("/upload", func(e echo.Context) error {
		return uploads.UploadFileHandler(e, store)
	})
	e.GET("/serve/", func(e echo.Context) error {
		return uploads.ServeFileHandler(e, store)
	})

	e.POST("/auth/login", lg.Login)
	e.POST("/auth/register", lg.Register)

	// Prepare upload directory
	err = os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		fmt.Println("Error creating the uploads directory:", err.Error())
		return
	}

	// Start Webserver
	if err = http.ListenAndServe(":8090", e); err != nil {
		fmt.Println("Error starting the webserver:", err.Error())
		return
	}
}

func pingHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Pong")
}
