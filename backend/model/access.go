package model

import "encoding/json"

/*访问表*/
type Access struct {
	ArticleID uint `gorm:"primaryKey"`
	AccessNum uint `gorm:"not null"`
}

func (a *Access) TableName() string {
	return "access"
}

func (a *Access) MarshalBinary() ([]byte, error) {
	return json.Marshal(a)
}

func (a *Access) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, a)
}
