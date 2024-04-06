package model

type Rft struct {
	Id            int    `gorm:"id"`
	Username      string `gorm:"username"`
	Refresh_token string `gorm:"refresh_token"`
}
