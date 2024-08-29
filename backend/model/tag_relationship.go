package model

import "encoding/json"

/*标签关系表*/
type TagRelationship struct {
	Name      string `gorm:"type:varchar(255);primary_key"`
	ArticleId uint   `gorm:"type:int;primary_key"`
}

func (tr *TagRelationship) TableName() string {
	return "tag_relationship"
}

func (tr *TagRelationship) MarshalBinary() ([]byte, error) {
	return json.Marshal(tr)
}
func (tr *TagRelationship) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, tr)
}
