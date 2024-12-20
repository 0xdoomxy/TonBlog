package model

import "time"

type AirportRelationship struct {
	AirportId   uint      `json:"airport_id" gorm:"primary_key"`
	UserAddress string    `gorm:"type:varchar(255);primary_key" json:"user_address"`
	CreateTime  time.Time `gorm:"type:datetime;not null" json:"create_time"`
	UpdateTime  time.Time `gorm:"type:datetime" json:"update_time"`
	FinishTime  time.Time `gorm:"type:datetime" json:"finish_time"`
	DeleteTime  time.Time `gorm:"type:datetime;" json:"delete_time"`
	Balance     float64   `gorm:"type:float"`
}

func (ar *AirportRelationship) TableName() string {
	return "airport_relationship"
}
