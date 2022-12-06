package main

import (
	"log"
	"os"

	"github.com/jaszczw/fiber/pkg/horde"
	"github.com/jaszczw/fiber/pkg/redis"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":3002"
	} else {
		port = ":" + port
	}

	return port
}

func main() { // Load the contents of the .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	redis.InitRedisClient()
	app := fiber.New()

	go redis.ListenRenderInRedis(func(requestId string) {
		horde.CheckImageStatusLoop(requestId)
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello, Railway!",
		})
	})

	app.Listen(getPort())
}
