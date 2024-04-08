package service

import (
	"log"
	"novel/dao"
	"novel/model"
)

func FindNovel(bookName string) (model.Novel, error) {
	novel, err := dao.SelectNovelByBookName(bookName)
	if err != nil {
		log.Println(err)
		return novel, err
	}
	return novel, nil
}
