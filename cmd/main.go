package main

import (
	"context"
	"fmt"
	"github.com/yayayapluto/revisi_api_lelang_online/cmd/config"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/utils"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("shutting down...")
		cancel()
	}()

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

	<-ctx.Done()
}
