package cmd

import (
	"fmt"

	"github.com/vviveksharma/lorebound/internal/database"
	"github.com/vviveksharma/lorebound/internal/logger"
	"github.com/vviveksharma/lorebound/internal/queue"
	"github.com/vviveksharma/lorebound/internal/server"
	"github.com/vviveksharma/lorebound/internal/utils"
)

func Execute() {
	logger.Init()

	db := database.Connect()
	if db == nil {
		logger.Fatal("database connection failed")
	}

	producer, err := queue.Connect()
	if err != nil {
		logger.Fatalf("error while connecting to the queue: %v", err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			logger.Error("kafka producer close", "err", err)
		}
	}()

	app := server.New(server.Deps{
		DB:       db,
		Producer: producer,
	})

	port := utils.Getenv("APP_PORT", "8080")
	addr := fmt.Sprintf(":%s", port)
	logger.Info("starting HTTP server", "addr", addr)
	if err := app.Listen(addr); err != nil {
		logger.Fatalf("error while starting server: %v", err)
	}
}
