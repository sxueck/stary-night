package storage

import (
	"fmt"
	"gorm.io/gorm"
	"lightning/config"
	"log"
	"reflect"
	"time"

	"gorm.io/driver/sqlite"
)

const DBTableName = "sites"

type DescribeSitesInfo struct {
	Name    string    `json:"name"`
	URL     string    `json:"url"`
	Author  string    `json:"author"`
	Lastmod time.Time `json:"lastmod"`
}

func ReSessionStorageConn() func() *gorm.DB {
	var db *gorm.DB

	return func() *gorm.DB {
		var err error
		if db == nil || reflect.DeepEqual(db, &gorm.DB{}) {
			db, err = gorm.Open(sqlite.Open(config.Cfg.DBName), &gorm.Config{
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

func LoadSitesToMemory(db func() *gorm.DB, memory chan<- []DescribeSitesInfo) error {
	var sitesCount int64 = 0
	sitesConn := db().Debug().Table(DBTableName)
	sitesConn.Count(&sitesCount)

	var ds = make([]DescribeSitesInfo, sitesCount)
	if sitesCount != 0 {
		err := sitesConn.Scan(&ds).Error
		if err != nil {
			return fmt.Errorf("The database fails to extract website information : %s\n", err)
		}
		memory <- ds
	}

	return nil
}

func AddMembers(db func() *gorm.DB, member DescribeSitesInfo) error {
	if reflect.DeepEqual(member, &DescribeSitesInfo{}) {
		return fmt.Errorf("please do not pass in empty members")
	}

	return db().Debug().Table(DBTableName).Create(&member).Error
}
