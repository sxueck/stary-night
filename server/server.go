package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"lightning/config"
	"lightning/storage"
	"log"
	"net/http"
)

type CustomContext struct {
	echo.Context // encapsulate the original context
}

func StartServ(ctx context.Context) {
	var errChan chan error
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.CORS())

	gpollers := GlobalPollers(ctx)
	webCtx := context.WithValue(ctx, "ds", gpollers)

	// use middleware to pass the context to handler
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &CustomContext{c}
			cc.Set("DescribeSites", webCtx)
			return next(cc)
		}
	})

	e.Static("/", "public")
	e.GET("/api/v1/list", ListAllSites)

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

func ListAllSites(c echo.Context) error {
	dss := c.Get("DescribeSites").(context.Context).
		Value("ds").(func() []storage.DescribeSitesInfo)

	fmt.Println(dss())

	sitesListResult, err := json.Marshal(dss())
	if err != nil {
		return c.String(http.StatusOK, err.Error())
	}

	return c.String(http.StatusOK, string(sitesListResult))
}
