package login

import "github.com/labstack/echo/v4"

type Login struct {
	d database
}

func (l Login) Register(c echo.Context) error {
	return nil
}

func (l Login) Login(c echo.Context) error {
	return nil
}

func (l Login) Setup() error {
	return l.d.create()
}
