package main

import (
	"github.com/nhatthanh123bk/E-commerce-website/db"
	"github.com/nhatthanh123bk/E-commerce-website/helper"
	"github.com/nhatthanh123bk/E-commerce-website/route"
)

func main() {
	helper.NewLogger()
	db.Init()
	db.InitRedis()
	e := route.Init()

	e.Logger.Fatal(e.Start(":1323"))
}
