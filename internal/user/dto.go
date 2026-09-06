package user

import "time"

type RegisterRequest struct {
	Nama         string `json:"nama" validate:"required,min=3"`
	Email        string `json:"email" validate:"required,email"`
	KataSandi    string `json:"kata_sandi" validate:"required,min=6"`
	NoTelp       string `json:"notelp" validate:"required"`
	TanggalLahir string `json:"tanggal_lahir" validate:"required"` // Format: YYYY-MM-DD
	JenisKelamin string `json:"jenis_kelamin" validate:"required,oneof=Laki-laki Perempuan"`
	Tentang      string `json:"tentang"`
	Pekerjaan    string `json:"pekerjaan"`
	IdProvinsi   string `json:"id_provinsi" validate:"required"`
	IdKota       string `json:"id_kota" validate:"required"`
}

type LoginRequest struct {
	Email     string `json:"email" validate:"required,email"`
	KataSandi string `json:"kata_sandi" validate:"required"`
}

// UserResponse hides the password and is used for returning user data
type UserResponse struct {
	ID           string    `json:"id"`
	Nama         string    `json:"nama"`
	Email        string    `json:"email"`
	NoTelp       string    `json:"notelp"`
	TanggalLahir time.Time `json:"tanggal_lahir"`
	JenisKelamin string    `json:"jenis_kelamin"`
	Tentang      string    `json:"tentang"`
	Pekerjaan    string    `json:"pekerjaan"`
	IdProvinsi   string    `json:"id_provinsi"`
	IdKota       string    `json:"id_kota"`
	IsAdmin      bool      `json:"is_admin"`
	CreatedAt    time.Time `json:"created_at"`
}

type UpdateProfileRequest struct {
	Nama         string `json:"nama" validate:"omitempty,min=3"`
	NoTelp       string `json:"notelp"`
	TanggalLahir string `json:"tanggal_lahir"` // Format: YYYY-MM-DD
	JenisKelamin string `json:"jenis_kelamin" validate:"omitempty,oneof=Laki-laki Perempuan"`
	Tentang      string `json:"tentang"`
	Pekerjaan    string `json:"pekerjaan"`
	IdProvinsi   string `json:"id_provinsi"`
	IdKota       string `json:"id_kota"`
}
