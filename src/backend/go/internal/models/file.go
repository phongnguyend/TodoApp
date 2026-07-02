package models

import "time"

// File is the GORM entity for uploaded-file metadata - analogous to an EF Core
// entity class. The on-disk `Location` is never exposed to clients directly;
// content is retrieved via the dedicated download endpoint instead.
type File struct {
	ID          uint       `gorm:"primarykey;autoIncrement"                json:"id"`
	Name        string     `gorm:"type:varchar(255);not null"              json:"name"`
	Extension   string     `gorm:"type:varchar(20);not null"               json:"extension"`
	Size        int64      `gorm:"not null"                                json:"size"`
	ContentType *string    `gorm:"type:varchar(100);column:content_type"   json:"contentType"`
	Location    string     `gorm:"type:varchar(500);not null"              json:"-"`
	CreatedAt   time.Time  `gorm:"autoCreateTime"                          json:"createdAt"`
	UpdatedAt   *time.Time `gorm:"autoUpdateTime"                          json:"updatedAt"`
}

// TableName overrides the default table name - mirrors [Table("files")] in EF.
func (File) TableName() string {
	return "files"
}
