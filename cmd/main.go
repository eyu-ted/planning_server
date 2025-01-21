package main

import (
	"time"

	"plan/config"
	route "plan/delivery/route"

	"github.com/gin-gonic/gin"
)

func main() {

	app := config.App()
	env := app.Env
	db := app.Mongo.Database(env.DBName)
	defer app.CloseDBConnection()


	timeout := time.Duration(env.ContextTimeout) * time.Second
	gin := gin.Default()

	route.Setup(env, timeout, db, gin)

	gin.Run()
}


