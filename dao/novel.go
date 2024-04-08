package dao

import (
	"fmt"
	"novel/model"
)

func SelectNovelByBookName(bookName string) (model.Novel, error) {
	var novel model.Novel
	err := dB.Table("novel").Where("book_name=?", bookName).Find(&novel)
	fmt.Println(err.Error)
	if err != nil {
		return novel, err.Error
	}
	return novel, nil
}
