package model

import "encoding/json"

/*文章关注总览表*/
type Like struct {
	ArticleID uint `gorm:"not null;Index:searchLike"`
	LikeNum   uint `gorm:"not null"`
}

func (like *Like) TableName() string {
	return "like"
}
func (like *Like) MarshalBinary() ([]byte, error) {
	return json.Marshal(like)
}

func (like *Like) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, like)
}
