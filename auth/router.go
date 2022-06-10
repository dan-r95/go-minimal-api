package auth

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

type Login struct {
	d      database
	Logger *zap.Logger
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegistrationRequest struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	Confirmation string `json:"confirmation"`
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
	db := &database{nil, l.Logger.Sugar()}
	l.d = *db
	if dbCon, err := l.d.setup(); err != nil {
		return err
	} else {
		db.db = dbCon
	}
	return nil
}

func present(c echo.Context, status int, pl interface{}, err error) error {
	if err == nil {
		return c.JSON(status, pl)
	} else {
		return c.JSON(status, echo.Map{"error": err.Error()})
	}
}
