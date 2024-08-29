package model

import (
	"encoding/json"
	"time"
)

/*评论表*/
type Comment struct {
	ID        uint `gorm:"primaryKey;autoIncrement"`
	CreateAt  time.Time
	TopID     uint   `gorm:"not null;index:search"`
	Content   string `gorm:"type:varchar(255);not null"`
	ArticleID uint   `gorm:"not null;index:search"`
	Creator   string `gorm:"varchar(64) not null;index:search"`
}

func (comment *Comment) TableName() string {
	return "comment"
}
func (comment *Comment) MarshalBinary() ([]byte, error) {
	return json.Marshal(comment)
}

func (comment *Comment) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, comment)
}
