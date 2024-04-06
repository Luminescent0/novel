package main

import (
	"novel/api"
	"novel/dao"
	"novel/service"
)

func main() {
	dao.InitDB()
	service.InitRdb()
	api.InitEngine()
}
