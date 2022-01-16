package server

import (
	"context"
	"lightning/spider"
	"lightning/storage"
	"log"
	"sync"
	"time"
)

func GlobalPollers(ctx context.Context) func() []storage.DescribeSitesInfo {
	var ds = make(chan []storage.DescribeSitesInfo, 1)
	db := storage.ReSessionStorageConn()
	var ticker = time.NewTicker(time.Second)
	var lowSLA []string

	go func(db func() *storage.DBConn, ds chan<- []storage.DescribeSitesInfo) {
		//	proactively update maintained memory information every hour
		for {
			select {
			case <-ticker.C:
				go func() {
					err := storage.LoadSitesToMemory(db(), ds)
					if err != nil {
						log.Println(err)
						return
					}
				}()
				ticker.Reset(30 * time.Minute)
			case <-ctx.Done():
				log.Println("exit the poller")
				return
			}
		}
	}(db, ds)

	var onlineSites []storage.DescribeSitesInfo
	go func() {
		for {
			select {
			case dss := <-ds:
				CheckTheSLA(dss, lowSLA)
				onlineSites = []storage.DescribeSitesInfo{}

				for _, v := range dss {
					var repetition = false
					for _, slaV := range lowSLA {
						if v.URL == slaV {
							repetition = !repetition
							break
						}
					}
					if !repetition {
						onlineSites = append(onlineSites, v)
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return func() []storage.DescribeSitesInfo {
		return onlineSites
	}
}

func CheckTheSLA(ds []storage.DescribeSitesInfo, lowSLA []string) {
	var wg sync.WaitGroup
	// this represents how many `goroutine`
	// are enabled to check for site survival
	ch := make(chan struct{}, 8)
	for i := 0; i < len(ds); i++ {
		ch <- struct{}{}
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			u := ds[i].URL
			isOnline := spider.SurvivalChecks(u)
			if !isOnline {
				lowSLA = append(lowSLA, u)
			}
			<-ch
		}(i)
	}
	wg.Wait()
}
