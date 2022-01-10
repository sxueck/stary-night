package server

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"lightning/config"
	"log"
)

func StartServ(ctx context.Context) {
	var errChan chan error
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.CORS())

	e.Static("/", "public")

	go func() {
		errChan <- e.Start(fmt.Sprintf(":%s", config.Cfg.Port))
	}()

	for {
		select {
		case err := <-errChan:
			log.Printf("an error occurred on the http server : %s\n", err)
			return
		case <-ctx.Done():
			log.Printf("the http service is exiting procedure\n")
			return
		}
	}
}
