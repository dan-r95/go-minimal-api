package main

import (
	"backend-homecase/login"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"strings"
)

var (
	supportedTypes = []string{"application/jpeg", "application/tiff", "application/png"}
)

func main() {
	// init logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	sugar := logger.Sugar()
	sugar.Infow("started my image api")

	lg := login.Login{}
	lg.Setup()

	e := echo.New()
	e.GET("/ping", pingHandler)
	e.POST("/upload", uploadFileHandler)
	e.GET("/serve/", serveFileHandler)

	e.POST("/auth/login", lg.Login)
	e.POST("/auth/register", lg.Register)

	// Prepare upload directory
	err := os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		fmt.Println("Error creating the uploads directory:", err.Error())
		return
	}

	// Start Webserver
	if err := http.ListenAndServe(":8090", e); err != nil {
		fmt.Println("Error starting the webserver:", err.Error())
		return
	}
}

func pingHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Pong")
}

func uploadFileHandler(c echo.Context) error {
	req := c.Request()
	// validate len and mime content type
	contentType := req.Header.Get("Content-type")
	if c.Request().ContentLength > 2*10^7 || !strings.Contains(strings.Join(supportedTypes, ","), contentType) {
		return c.String(http.StatusBadRequest, "Supported file types: "+strings.Join(supportedTypes, ","))

	}
	file, fileHandler, err := c.Request().FormFile("file")
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error when parsing the uploaded file: "+err.Error())
	}

	// add auth middle ware

	defer file.Close()

	// Create a new file in the uploads directory
	destFile, err := os.Create(fmt.Sprintf("./uploads/%s", fileHandler.Filename))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error when trying to write the uploaded file to the upload directory:"+err.Error())
	}

	defer destFile.Close()

	// Write uploaded file to filesystem
	_, err = io.Copy(destFile, file)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error when trying to write the uploaded file to the upload directory: "+err.Error())
	}

	return c.String(http.StatusOK, "File successfully uploaded!")
}

func serveFileHandler(c echo.Context) error {
	req := c.Request()
	hasSize(req)
	requestedFile := strings.TrimPrefix(req.URL.Path, "/serve/")
	return c.File("uploads/" + requestedFile)
}

func hasSize(req *http.Request) string {
	params := req.URL.Query()
	return params.Get("size")
}
