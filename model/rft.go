package model

type Rft struct {
	Id           int
	Username     string
	RefreshToken string `gorm:"column:refresh_token"`
}
