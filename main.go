package main

import (
	"github.com/blogs/db"
	"github.com/blogs/helper"
	"github.com/blogs/route"
)

func main() {
	helper.NewLogger()
	db.Init()
	db.InitRedis()
	e := route.Init()

	e.Logger.Fatal(e.Start(":1323"))
}
