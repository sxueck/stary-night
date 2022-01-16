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
	Name   string `json:"name"`
	URL    string `json:"url"`
	Author string `json:"author"`

	// Can use QQ / weixin / mailbox three types
	// acceptance is not differentiated
	Contact     string    `json:"contact"`
	Lastmod     time.Time `json:"lastmod"`
	Description string    `json:"description"`
}

type DBConn struct {
	*gorm.DB
}

func ReSessionStorageConn() func() *DBConn {
	var db *gorm.DB

	return func() *DBConn {
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
		return &DBConn{db}
	}
}

func LoadSitesToMemory(db *DBConn, memory chan<- []DescribeSitesInfo) error {
	var sitesCount int64 = 0
	sitesConn := db.Debug().Table(DBTableName)
	sitesConn.Count(&sitesCount)

	var ds = make([]DescribeSitesInfo, sitesCount)
	if sitesCount != 0 {
		err := sitesConn.Scan(&ds).Error
		if err != nil {
			return fmt.Errorf("The database fails to extract website information : %s\n", err)
		}
		log.Println("data cache update completed")
		memory <- ds
	}

	return nil
}

func (d *DBConn) AddMembers(member DescribeSitesInfo) error {
	if reflect.DeepEqual(member, &DescribeSitesInfo{}) {
		return fmt.Errorf("please do not pass in empty members")
	}

	if haveRepetition, err := d.repeatedSiteChecks(member); err != nil {
		return fmt.Errorf("errors occurred while checking for duplicates ：%s", err)
	} else {
		if haveRepetition {
			return fmt.Errorf("this url already exists")
		}
	}

	return d.Debug().Table(DBTableName).Create(&member).Error
}

func (d *DBConn) repeatedSiteChecks(member DescribeSitesInfo) (bool, error) {
	var count int64 = 0
	if err := d.Debug().
		Table(DBTableName).
		Where("url = '?'", member.URL).
		Count(&count).Error; err != nil {
		return false, err
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
}
