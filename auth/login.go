package auth

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Login struct {
	d      database
	Logger *zap.Logger
}

func (l Login) Register(c echo.Context) error {
	return nil
}

func (l Login) Login(c echo.Context) error {
	return nil
}

func (l Login) Setup() error {
	l.d = database{nil, l.Logger.Sugar()}
	return l.d.create()
}
