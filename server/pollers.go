package server

import (
	"context"
	"gorm.io/gorm"
	"lightning/storage"
	"log"
	"time"
)

func GlobalPollers(ctx context.Context) func() []storage.DescribeSitesInfo {
	var ds = make(chan []storage.DescribeSitesInfo, 1)
	db := storage.ReSessionStorageConn()
	var ticker = time.NewTicker(time.Hour)

	go func(db func() *gorm.DB, ds chan<- []storage.DescribeSitesInfo) {
		go func() {
			err := storage.LoadSitesToMemory(db, ds)
			if err != nil {
				log.Println(err)
				return
			}
		}()

		//	proactively update maintained memory information every hour
		for {
			select {
			case <-ticker.C:
				ticker.Reset(time.Hour)
			case <-ctx.Done():
				log.Println("exit the poller")
				return
			}
		}
	}(db, ds)

	var dss []storage.DescribeSitesInfo
	go func() {
		for {
			select {
			case dss = <-ds:
			case <-ctx.Done():
				return
			}
		}
	}()

	return func() []storage.DescribeSitesInfo {
		return dss
	}
}
