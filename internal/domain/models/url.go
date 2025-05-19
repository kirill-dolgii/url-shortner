package models

type Url struct {
	ID      int64  `gorm:"primaryKey;column:id"`
	FullUrl string `gorm:"column:url"`
	Alias   string `gorm:"column:alias"`
}

func (Url) TableName() string {
	return "url"
}
