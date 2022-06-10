package main

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go-minimal-api/auth"
	"go-minimal-api/uploads"
	"go.uber.org/zap"
	"net/http"
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
		sugar.Errorw("Error setting up the db:", err.Error())
	}

	store, _ := uploads.NewStorage(context.Background())
	if err != nil {
		sugar.Errorw("Error setting up the storage:", err.Error())
	}

	e := echo.New()
	// login without authentication
	g := e.Group("/auth")
	g.POST("/login", lg.Login)

	g2 := e.Group("/")
	// add token middleware (login should be without token
	g2.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:  []byte("secret"),
		TokenLookup: "query:token",
	}))
	g2.POST("/register", lg.Register)
	g2.GET("/ping", pingHandler)
	g2.POST("/upload", func(e echo.Context) error {
		return uploads.UploadFileHandler(e, store)
	})
	e.GET("/serve/", func(e echo.Context) error {
		return uploads.ServeFileHandler(e, store)
	})

	// Start Webserver
	if err = http.ListenAndServe(":8090", e); err != nil {
		sugar.Errorw("Error starting the webserver:", err.Error())
		return
	}
}

func pingHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Pong")
}
