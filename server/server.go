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
	"net/http"
	"reflect"
)

/*
200: Default
400: Method exception or incorrectly formatted data submitted
403: Abnormal access
404: Resources( page / db record ...) do not exist
502: External component is abnormal procedure
*/

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

	e.GET("/api/v1/list", ListAllSites)
	e.POST("/api/v1/site", AddMembersHandler)
	e.GET("/api/v1/ran_url", RandomSite)
	e.POST("/api/v1/subscribe", SubscribeUpdate)
	e.Static("/", "public")

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
	// debug
	//body, _ := ioutil.ReadAll(c.Request().Body)
	//log.Println(string(body))
	//c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(body))

	cc := c.(*CustomContext)
	var ds = storage.DescribeSitesInfo{}
	err := c.Bind(&ds)
	if err != nil {
		return c.String(http.StatusOK,
			fmt.Sprintf("[ERROR] error parsing user json : %s", err))
	}

	if len(ds.URL) == 0 {
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
		return c.String(http.StatusOK,
			fmt.Sprintf("[ERROR] error writing to database : %s", err))
	}
	return c.String(http.StatusOK, ds.Name)
}

func RandomSite(c echo.Context) error {
	dss := c.Get("DescribeSites").(context.Context).
		Value("ds").(func() []storage.DescribeSitesInfo)

	if len(dss()) == 0 {
		return c.String(http.StatusOK, "[ERROR] currently no web site exists")
	}

	token, err := NewUserAccess(c)
	if err != nil {
		return c.String(http.StatusBadGateway, storage.ReDBHaveError(err))
	}

	v, ok := globalServerSessionCache[token]
	if !ok {
		return c.String(http.StatusOK,
			"An unknown error occurred and the user cache does not exist")
	}

	result, err := json.Marshal(dss()[v.Card[v.Pos]])
	if err != nil {
		return c.String(http.StatusOK,
			fmt.Sprintf("[ERROR] gives random site errors : %s", err))
	}

	return c.String(http.StatusOK, string(result))
}

func SubscribeUpdate(c echo.Context) error {
	cc := c.(*CustomContext)
	var subs = &storage.SubscribeMembers{}
	err := c.Bind(subs)
	if err != nil {
		return c.String(http.StatusBadRequest,
			fmt.Sprintf("there is an anomaly in the json you submitted : %s", err))
	}

	if reflect.DeepEqual(subs, &storage.SubscribeMembers{}) {
		return c.String(http.StatusBadRequest,
			fmt.Sprintf("please check the data you submitted\n"))
	}

	db := cc.GetDBConn()
	if found, err := db().SelectSubscribe(subs.Mail); err != nil {
		return c.String(http.StatusBadGateway, storage.ReDBHaveError(err))
	} else {
		if found != 0 {
			return c.String(http.StatusOK, "the record already exists")
		}
	}

	err = db().AddSubscribeRoll(*subs)
	if err != nil {
		return c.String(http.StatusOK,
			fmt.Sprintf("an error occurred while adding subscriber : %s\n", err))
	}

	return c.String(http.StatusOK, "Success")
}
