package dao

import (
	"fmt"
	"log"
	"novel/model"
)

func SelectUserByUsername(username string) (model.User, error) {
	var user model.User
	err := dB.Table("user").Where("username=?", username).Find(&user)
	fmt.Println(err.Error)
	if err != nil {
		return user, err.Error
	}
	return user, nil
}
func InsertUser(user model.User) error {
	err := dB.Table("user").Select("username", "password").Create(&user)
	if err != nil {
		fmt.Println(err.Error)
		return err.Error
	}
	return nil
}

func UpdatePassword(username, newPassword string) error {
	user := model.User{Username: username, Password: newPassword}
	err := dB.Table("user").Model(&user).Where("username=?", username).Update("password", newPassword)
	if err != nil {
		log.Println(err.Error)
		return err.Error
	}
	return nil
}
