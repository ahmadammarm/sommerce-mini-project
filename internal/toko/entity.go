package toko

import (
	"time"

	"github.com/lucsky/cuid"
	"gorm.io/gorm"
)

type Toko struct {
	ID        string    `gorm:"primaryKey;type:varchar(50);column:id"`
	IdUser    string    `gorm:"column:id_user;type:varchar(50);not null"`
	NamaToko  string    `gorm:"column:nama_toko;type:varchar(255)"`
	UrlFoto   string    `gorm:"column:url_foto;type:varchar(255)"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (t *Toko) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = cuid.New()
	return
}
