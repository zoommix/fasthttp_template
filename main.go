package main

import (
	"github.com/alexsasharegan/dotenv"
	"github.com/valyala/fasthttp"
	"github.com/zoommix/fasthttp_template/router"
	pg "github.com/zoommix/fasthttp_template/store"
	"github.com/zoommix/fasthttp_template/utils"
)

func main() {
	defer pg.Close()

	err := dotenv.Load()
	utils.InitLogger()

	if err != nil {
		utils.LogError("Error loading .env file: ", err.Error())
	}

	r := router.New()
	p := utils.GetPort()

	utils.LogNotice("Starting server...")
	utils.LogNotice("Listening on tcp://0.0.0.0:" + p)

	err = fasthttp.ListenAndServe(":"+p, r.Handler)

	if err != nil {
		panic(err)
	}
}
