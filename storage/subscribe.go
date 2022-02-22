package storage

import "time"

// the number of members used to subscribe

const DBSubsTableName = "sites"

type SubscribeMembers struct {
	Mail       string    `json:"mail"`
	Insert     time.Time `json:"insert"`
	WebManager bool      `json:"webmanager"`
}

func (d *DBConn) AddSubscribeRoll(members SubscribeMembers) {

}
