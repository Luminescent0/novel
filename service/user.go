package service

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"novel/dao"
	"novel/model"
	"time"
)

func IsPasswordCorrect(username, password string) (bool, error) {
	user, err := dao.SelectUserByUsername(username)
	if err != nil {
		fmt.Println(err)
		flag := errors.Is(err, gorm.ErrRecordNotFound)
		if !flag {
			return false, err
		}
		fmt.Println(username)
		return false, err
	}
	flag := ComparePassword(user.Password, []byte(password))
	if !flag {
		return false, nil
	}
	fmt.Println("验证密码成功")
	return true, nil
}
func ComparePassword(hashedPassword string, plainPassword []byte) bool {
	byteHash := []byte(hashedPassword)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPassword)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
func UsernameIsExist(username string) error {
	_, err := dao.SelectUserByUsername(username)
	if err != nil {
		fmt.Println(err)
		if err == gorm.ErrRecordNotFound {
			fmt.Println("用户名不存在")
			return err
		}
		return err
	}
	return nil
}
func IsRepeatUsername(username string) (bool, error) {
	_, err := dao.SelectUserByUsername(username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func Register(user model.User) error {
	password, err := Cipher(user)
	if err != nil {
		return err
	}
	user.Password = password
	err = dao.InsertUser(user)
	if err != nil {
		return err
	}
	return nil
}

func Cipher(user model.User) (string, error) {
	password := []byte(user.Password)
	nowG := time.Now()
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	if err != nil {
		fmt.Println("加密后:", string(hashedPassword), "耗时", time.Now().Sub(nowG))
	}
	return string(hashedPassword), nil
}

func ChangePassword(username, newPassword string) error {
	user := model.User{Username: username, Password: newPassword}
	cnewPassword, err := Cipher(user)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = dao.UpdatePassword(username, cnewPassword)
	fmt.Println(err)
	return err
}

func SelectUserByUsername(username string) (model.User, error) {
	user, err := dao.SelectUserByUsername(username)
	if err != nil {
		log.Println(err)
		if err == gorm.ErrRecordNotFound {
			fmt.Println("用户不存在")
			return user, err
		}
		return user, err
	}
	return user, nil
}

func SelectLiked(bookId, userId int) (bool, error) {
	err := dao.SelectLiked(bookId, userId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func CancelLiked(bookId, userId int) error {
	err := dao.DelLiked(bookId, userId)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func Liked(bookId, userId int) error {
	err := dao.AddLiked(bookId, userId)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
