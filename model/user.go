package model

type User struct {
	Id           int    `gorm:"primaryKey"`
	Username     string `json:"username" validate:"min=4,max=10"`
	Password     string `json:"password" validate:"min=6,max=16"`
	RefreshToken string `json:"refresh_token"`
}
