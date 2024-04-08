package model

type User struct {
	Id       int    `gorm:"PRIMARY_KEY"`
	Username string `json:"username" validate:"min=4,max=10"`
	Password string `json:"password" validate:"min=6,max=16"`
}

type Liked struct {
	Id     int `gorm:"PRIMARY_KEY"`
	BookId int
	UserId int
}
