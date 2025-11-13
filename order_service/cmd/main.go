package main

import (
	ordersApi "github.com/ilovepitsa/orders/api/order"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func main() {
	e := echo.New()
	sh := ordersApi.NewStrictHandler(nil, nil)
	ordersApi.RegisterHandlers(e, sh)

	err := e.Start("0.0.0.0:8071")
	log.Error(err)
}
