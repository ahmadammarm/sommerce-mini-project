package category

import (
	"time"

	"github.com/lucsky/cuid"
	"gorm.io/gorm"
)

type Category struct {
	ID           string    `gorm:"primaryKey;type:varchar(50);column:id"`
	NamaCategory string    `gorm:"column:nama_category;type:varchar(255)"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (c *Category) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = cuid.New()
	return
}
