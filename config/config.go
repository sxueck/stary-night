package config

import (
	"fmt"
	"log"
	"os"
	"reflect"
)

var Cfg = &cfg{}

type cfg struct {
	Port string `env:"HTTP_PORT" envDefault:"80"`
}

func ArgsEnv(st interface{}) {
	rVal := reflect.ValueOf(st)
	rType := reflect.TypeOf(st)
	if rType.Kind() != reflect.Ptr {
		log.Fatalln("please pass with pointer type, like &args")
	} else {
		rVal = rVal.Elem()
		rType = rType.Elem()
	}

	for i := 0; i < rVal.NumField(); i++ {
		t := rType.Field(i)
		f, _ := rType.FieldByName(t.Name)

		env := f.Tag.Get("env")
		if len(env) == 0 {
			continue
		}

		eArgs := os.Getenv(env)
		if len(eArgs) == 0 {
			eArgs = f.Tag.Get("envDefault")
		}

		if fn := rVal.FieldByName(t.Name); fn.CanSet() {
			fn.SetString(eArgs)
		}

	}
}

func init() {
	ArgsEnv(Cfg)
	fmt.Printf("%+v\n", Cfg)
}
