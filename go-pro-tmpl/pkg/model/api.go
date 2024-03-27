package model

import "time"

type APIRegistryInfo struct {
	ID       string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	Type     string    `gorm:"type:varchar(20)" json:"type"`
	Content  string    `gorm:"type:longtext" json:"content"`
	From     string    `gorm:"column:from;type:varchar(20)" json:"from"`
	CreateAt time.Time `gorm:"type:datetime" json:"create_at"`
	UpdateAt time.Time `gorm:"type:datetime" json:"update_at"`
}

// TableName This method is used to specify the table name for GORM
func (APIRegistryInfo) TableName() string {

	return "api_registry_info"
}
