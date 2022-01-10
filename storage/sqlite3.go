package storage

import (
	"gorm.io/gorm"
	"log"
	"reflect"

	"gorm.io/driver/sqlite"
)

type DescribeSitesInfo struct {
}

func ReSessionStorageConn() func() *gorm.DB {
	var db *gorm.DB

	return func() *gorm.DB {
		var err error
		if db == nil || reflect.DeepEqual(db, &gorm.DB{}) {
			db, err = gorm.Open(sqlite.Open("session.db"), &gorm.Config{
				QueryFields: false,
			})
			if err != nil {
				log.Println("unable to connect to database")
				return nil
			}
		}
		return db
	}
}

func LoadSitesToMemory() {

}
