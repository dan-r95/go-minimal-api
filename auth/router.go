package auth

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

type (
	LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	RegistrationRequest struct {
		Email        string `json:"email"`
		Password     string `json:"password"`
		Confirmation string `json:"confirmation"`
	}

	Login struct {
		d      Database
		Logger *zap.Logger
	}
)

func NewLogin(logger *zap.SugaredLogger) (Database, error) {
	return newDatabase(logger)
}

func newLogin(logger *zap.SugaredLogger) (Database, error) {
}

func (l Login) Register(c echo.Context) (err error) {
	var req RegistrationRequest
	if err = c.Bind(&req); err == nil {
		err = l.d.Register(req)
		return present(c, http.StatusCreated, nil, err)
	}
	return
}

func (l Login) Login(c echo.Context) (err error) {
	var req LoginRequest
	if err = c.Bind(&req); err == nil {
		err = l.d.Login(req)
		return present(c, http.StatusCreated, nil, err)
	}
	return
}

func (l Login) Setup() error {
	db, err := New(l.Logger.Sugar())
	if err != nil {
		return err
	}
	l.d = db
	return nil
}

func present(c echo.Context, status int, pl interface{}, err error) error {
	if err == nil {
		return c.JSON(status, pl)
	} else {
		return c.JSON(status, echo.Map{"error": err.Error()})
	}
}
