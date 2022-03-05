package storage

import "time"

// the number of members used to subscribe

const DBSubsTableName = "subscribes"

type SubscribeMembers struct {
	Mail       string    `json:"mail"`
	Insert     time.Time `json:"insert"`
	Webmanager bool      `json:"webmanager"`
}

func (d *DBConn) AddSubscribeRoll(members SubscribeMembers) error {
	if members.Insert.IsZero() {
		members.Insert = time.Now()
	}

	if err := d.Debug().
		Table(DBSubsTableName).
		Create(members).
		Error; err != nil {
		return err
	}

	return nil
}

func (d *DBConn) SelectSubscribe(mail string) (int64, error) {
	var count int64 = 0
	err := d.Debug().
		Table(DBSubsTableName).
		Where("mail = ?", mail).
		Count(&count).Error

	return count, err
}
