package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
	"lightning/config"
	"lightning/storage"
	"log"
	"net/http"
)

type CustomContext struct {
	echo.Context // encapsulate the original context
}

func (cc *CustomContext) GetDBConn() func() *gorm.DB {
	return nil
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
		errChan <- e.Start(fmt.Sprintf("%s:%s", config.Cfg.Address, config.Cfg.Port))
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

	sitesListResult, err := json.Marshal(dss())
	if err != nil {
		return c.String(http.StatusOK, err.Error())
	}

	return c.String(http.StatusOK, string(sitesListResult))
}

func AddMembersHandler(c echo.Context) error {
	var ds = storage.DescribeSitesInfo{}
	err := c.Bind(&ds)
	if err != nil {
		return c.String(http.StatusInternalServerError,
			fmt.Sprintf("error parsing user json : %s", err))
	}
	return nil
}
