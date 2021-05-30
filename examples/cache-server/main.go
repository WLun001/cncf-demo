package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// default to no store
	app.Use(func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderCacheControl, "no-store")
		return c.Next()
	})

	app.Get("/", func(c *fiber.Ctx) error {
		// simulate database call
		// for 200ms
		time.Sleep(200 * time.Millisecond)
		return c.JSON(fiber.Map{"result": "ok"})
	})

	app.Get("/cache", allowCache(5), func(c *fiber.Ctx) error {
		// simulate database call
		// for 200ms
		time.Sleep(200 * time.Millisecond)
		return c.JSON(fiber.Map{"result": "ok"})
	})
	log.Fatal(app.Listen(":3000"))
}

func allowCache(dur int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderCacheControl, fmt.Sprintf("public, max-age=%d", dur))
		return c.Next()
	}
}
