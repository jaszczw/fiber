package main

import (
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
	godotenv.Load(".env")

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
