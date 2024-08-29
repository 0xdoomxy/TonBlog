package model

import (
	"encoding/json"
	"time"
)

/*用户表*/
type User struct {
	Address   string `gorm:"type:varchar(64);primary_key"`
	Alias     string `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time
	//TODO we should pass this field to notify the user who starts up the subscribe function  the new message is arriving
	//TgAccount string `gorm:"type:varchar(255);not null"`
}

func (user *User) TableName() string {
	return "user"
}

func (user *User) MarshalBinary() ([]byte, error) {
	return json.Marshal(user)
}

func (user *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, user)
}
