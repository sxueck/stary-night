package main

import (
	"fmt"
	"github.com/emersion/go-sasl"
	"github.com/labstack/echo/v4"
	"stary-night/config"
)

type Message struct {
	UserName     string `json:"user_name"`
	SiteName     string `json:"site_name"`
	ReviewStatus string `json:"review_status"`
	Theme        string `json:"theme"`

	// other
	DefineMessage string `json:"define_message"`
	SendObject    string `json:"send_object"`
}

type cfg struct {
	Server        string `env:"SERVER"`
	Password      string `env:"SMTP_PASSWORD"`
	UserName      string `env:"SMTP_PASSWORD"`
	InterfacePort string `env:"SERVER_PORT" envDefault:"80"`
}

var Cfg = &cfg{}

func init() {
	config.ArgsEnv(&Cfg)
	fmt.Printf("%+v\n", Cfg)
}

func main() {
	auth := sasl.NewPlainClient("", Cfg.UserName, Cfg.Password)

}

func StartServ() {
	e := echo.New()
	e.HideBanner = true

	e.POST("/api/v1/send", func(c echo.Context) error {
		var
		return nil
	})
}
