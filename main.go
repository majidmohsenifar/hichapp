package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/majidmohsenifar/hichapp/cmd"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	cmd.RunHttpServer()
}
