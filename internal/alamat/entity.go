package alamat

import (
	"time"
)

type Alamat struct {
	ID           string    `gorm:"primaryKey;type:varchar(50);column:id"`
	IdUser       string    `gorm:"column:id_user;type:varchar(50);not null"`
	JudulAlamat  string    `gorm:"column:judul_alamat;type:varchar(255)"`
	NamaPenerima string    `gorm:"column:nama_penerima;type:varchar(255)"`
	NoTelp       string    `gorm:"column:no_telp;type:varchar(255)"`
	DetailAlamat string    `gorm:"column:detail_alamat;type:text"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
}
