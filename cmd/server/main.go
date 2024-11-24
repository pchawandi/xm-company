package main

import (
	"context"
	"log"

	"github.com/pchawandi/xm-company/api"
	"github.com/pchawandi/xm-company/database"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	db := database.NewDatabase()
	dbWrapper := &database.GormDatabase{DB: db}
	ctx := context.Background()
	logger, _ := zap.NewProduction()
	defer func() {
		_ = logger.Sync()
	}()

	//gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.DebugMode)

	r := api.NewRouter(ctx, dbWrapper, logger)

	if err := r.Run(":8001"); err != nil {
		log.Fatal(err)
	}
}
