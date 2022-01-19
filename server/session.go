package server

import (
	"github.com/labstack/echo/v4"
	"lightning/storage"
	"log"
	"math/rand"
	"sync"
	"time"
)

var globalServerSessionCache = make(map[string]SessionColumn)

// Adding IP judgment prevents spoofing
// of cache space by forging userAgent
//var ipCount = make(map[string]int)

type SessionColumn struct {
	Expired   time.Time
	Pos       int
	Card      []int
	MaxLength int
	Session   string
}

func NewUserAccess(c echo.Context) (string, error) {
	var mu sync.Mutex
	//ip := c.RealIP()
	//
	//if v, ok := ipCount[ip]; ok {
	//	// limit
	//	if v >= 5 {
	//		return "", fmt.Errorf("the access limit for a single ip address was exceeded")
	//	}
	//}
	//
	//ipCount[ip]++

	token := ComputeUserRepresentationId(c.RealIP(), c.Request().UserAgent())
	expired := time.Now().Add(10 * time.Minute)
	if v, ok := globalServerSessionCache[token]; ok {
		// update the cache
		v = SessionColumn{
			Expired:   expired,
			Pos:       v.Pos + 1,
			Card:      v.Card,
			MaxLength: v.MaxLength,
			Session:   token,
		}

		if v.Pos == v.MaxLength {
			v.Pos = 0
		}

		mu.Lock()
		globalServerSessionCache[token] = v
		mu.Unlock()
	} else {
		db := storage.ReSessionStorageConn()
		v = SessionColumn{
			Expired:   expired,
			Pos:       0,
			MaxLength: db().MembersCount(),
			Session:   token,
		}
		v.Card = rand.Perm(v.MaxLength)
		globalServerSessionCache[token] = v
	}

	log.Println(globalServerSessionCache)
	return token, nil
}
