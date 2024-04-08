package dao

import (
	"log"
	"novel/model"
)

func AddRefreshToken(username, token string) error {
	var rf model.Rft
	res := dB.Table("rft").Where("username=?", username).Find(&rf)
	if res.RowsAffected == 0 {
		rf.Username = username
		rf.RefreshToken = token
		err := dB.Table("rft").Create(&rf)
		if err != nil {
			log.Println(err)
			return err.Error
		}
		return nil
	}
	err := dB.Table("rft").Where("username=?", username).Update("refresh_token", token)
	return err.Error
}

func GetRefreshToken(rt string) string {
	var rf model.Rft
	dB.Table("rft").Where("refresh_token = ?", rt).Scan(&rf)
	return rf.Username
}

func DelRefreshToken(rt, username string) {
	err := dB.Table("rft").Where("refresh_token = ? AND username= ?", rt, username).Update("refresh_token", nil)
	if err != nil {
		log.Println(err)
	}
	return
}
