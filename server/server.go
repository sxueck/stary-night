package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
	"html/template"
	"io"
	"lightning/config"
	"lightning/storage"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type CustomContext struct {
	echo.Context // encapsulate the original context
}

type EchoTemplate struct {
	templates *template.Template
}

func (t *EchoTemplate) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func (cc *CustomContext) GetDBConn() func() *gorm.DB {
	return nil
}

func StartServ(ctx context.Context) {
	var errChan chan error
	e := echo.New()
	e.HideBanner = true
	e.Renderer = &EchoTemplate{}

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

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if found, _ := regexp.MatchString("\\.(html|htm)", c.Request().RequestURI); found {
				return c.String(http.StatusForbidden,
					"please do not access web resources directly")
			}

			return next(c)
		}
	})

	//e.Static("/", "public")
	e.Static("/static", "public/static")
	e.Any("/", RenderStaticPages)
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

func RenderStaticPages(c echo.Context) error {
	reqURL := c.Request().RequestURI
	var resNameSplitIndex = strings.LastIndex(reqURL, "/")
	var resName = reqURL[resNameSplitIndex:]
	if found, _ := regexp.MatchString(
		"[a-zA-z]+://[^\\s]*?/[^\\s]*?[/|.]", reqURL); !found {
		// like x.x/1.jpg
		resName = fmt.Sprintf("%s.html", resName)
	}

	return c.Render(http.StatusOK, resName, nil)
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
