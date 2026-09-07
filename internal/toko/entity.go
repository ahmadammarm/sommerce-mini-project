package toko

import (
	"time"
)

type Toko struct {
	ID        string    `gorm:"primaryKey;type:varchar(50);column:id"`
	IdUser    string    `gorm:"column:id_user;type:varchar(50);not null;uniqueIndex"`
	NamaToko  string    `gorm:"column:nama_toko;type:varchar(255)"`
	UrlFoto   string    `gorm:"column:url_foto;type:varchar(255)"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}
