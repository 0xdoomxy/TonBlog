package model

import "encoding/json"

/*文章关注表*/
type LikeRelationship struct {
	ArticleID uint   `gorm:"primarykey"`
	Address   string `gorm:"address;varchar(64)"`
}

func (lrs *LikeRelationship) TableName() string {
	return "like_relationship"
}
func (lrs *LikeRelationship) MarshalBinary() ([]byte, error) {
	return json.Marshal(lrs)
}

func (lrs *LikeRelationship) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, lrs)
}
