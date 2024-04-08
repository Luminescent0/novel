package model

import "time"

type Novel struct {
	Id         int
	BookName   string
	Brief      string
	Stat       int
	Typ        string `gorm:"column:type"`
	UpdateTime time.Time
}

type Chapter struct {
	Id          int `gorm:"PRIMARY_KEY"`
	ChapterName string
	Content     string
}
