package dao

import (
	"fmt"
	"log"
	"novel/model"
)

func SelectUserByUsername(username string) (model.User, error) {
	user := model.User{}
	err := dB.Table("user").Where("username=?", username).Take(&user)
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

func AddRefreshToken(username, token string) error {
	return dB.Create(&model.User{
		Username:     username,
		RefreshToken: token,
	}).Error
}

func GetRefreshToken(rt string) string {
	var username string
	dB.Model(&model.User{}).Select("username").
		Where("refresh_token = ? ", rt).
		Scan(&username)
	return username
}

func DelRefreshToken(rt, username string) {
	err := dB.Model(&model.User{}).Where("refresh_token = ? AND username= ?", rt, username).UpdateColumn("refresh_token", nil)
	if err != nil {
		log.Println(err)
	}
	return
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
