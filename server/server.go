package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"lightning/config"
	"lightning/spider"
	"lightning/storage"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type CustomContext struct {
	echo.Context // encapsulate the original context
}

func (cc *CustomContext) GetDBConn() func() *storage.DBConn {
	return storage.ReSessionStorageConn()
}

func StartServ(ctx context.Context) {
	var errChan chan error
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

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
	e.POST("/api/v1/site", AddMembersHandler)
	e.GET("/api/v1/ran_url", RandomSite)

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
	cc := c.(*CustomContext)
	var ds = storage.DescribeSitesInfo{}
	err := c.Bind(&ds)
	if err != nil {
		return c.String(http.StatusOK, fmt.Sprintf("[ERROR] error parsing user json : %s", err))
	}

	if ds.URL == "" {
		return c.String(http.StatusNotFound, "[ERROR] parameter error")
	}

	// here for future users to judge their
	// own data is correct to leave the interface
	if ds.Name == "" {
		ds.Name = spider.ObtainSiteTitle(ds.URL)
	}

	db := cc.GetDBConn()
	err = db().AddMembers(ds)
	if err != nil {
		return c.String(http.StatusOK, fmt.Sprintf("[ERROR] error writing to database : %s", err))
	}
	return c.String(http.StatusOK, ds.Name)
}

func RandomSite(c echo.Context) error {
	dss := c.Get("DescribeSites").(context.Context).
		Value("ds").(func() []storage.DescribeSitesInfo)

	rand.Seed(time.Now().Unix())
	ran := rand.Intn(len(dss()))

	result, err := json.Marshal(dss()[ran])
	if err != nil {
		return c.String(http.StatusOK, fmt.Sprintf("[ERROR] gives random site errors : %s", err))
	}

	return c.String(http.StatusOK, string(result))
}
