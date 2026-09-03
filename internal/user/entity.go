package user

import (
	"time"
)

type User struct {
	ID           string    `gorm:"primaryKey;type:varchar(50);column:id"`
	Nama         string    `gorm:"column:nama;type:varchar(255)"`
	KataSandi    string    `gorm:"column:kata_sandi;type:varchar(255)"`
	NoTelp       string    `gorm:"column:notelp;type:varchar(255);unique"`
	TanggalLahir time.Time `gorm:"column:tanggal_lahir;type:date"`
	JenisKelamin string    `gorm:"column:jenis_kelamin;type:varchar(255)"`
	Tentang      string    `gorm:"column:tentang;type:text"`
	Pekerjaan    string    `gorm:"column:pekerjaan;type:varchar(255)"`
	Email        string    `gorm:"column:email;type:varchar(255);unique"`
	IdProvinsi   string    `gorm:"column:id_provinsi;type:varchar(255)"`
	IdKota       string    `gorm:"column:id_kota;type:varchar(255)"`
	IsAdmin      bool      `gorm:"column:isAdmin;type:boolean;default:false"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
}
