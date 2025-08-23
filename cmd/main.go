package main

import (
	"fmt"
	"github.com/yayayapluto/revisi_api_lelang_online/cmd/config"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/utils"
)

func main() {
	env, err := utils.LoadEnv()
	fmt.Println(env)

	if err != nil {
		panic(err)
	}

	db, err := config.ConnectDB(env.DBHOST, env.DBUSER, env.DBPASSWORD, env.DBNAME, env.DBPORT)
	if err != nil {
		panic(err)
	}

	app, err := config.NewApp(db)
	if err != nil {
		panic(err)
	}

	if err := app.Listen(":8080"); err != nil {
		panic(err)
	}

	fmt.Println("Success connected to database!")
}
