package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
)

// FruitInterface is a collection of various fruits.
type FruitInterface struct {
	Page   int
	Fruits []string
}

func main() {
	app := fiber.New()

	// Load static files like CSS, Images & JavaScript.
	app.Static("/static", "./static")

	// Returns a local HTML file.
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("./templates/hello.html")
		// navigate to => http://localhost:3000/
	})

	// Returns plain text.
	app.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("Hello, Fiber!")
		// navigate to => http://localhost:3000/hello
	})

	// Use parameters.
	app.Get("/parameter/:value", func(c *fiber.Ctx) error {
		return c.SendString("Get request with value: " + c.Params("value"))
		// navigate to => http://localhost:3000/parameter/this_is_the_parameter
	})

	// Use wildcards to design your API.
	app.Get("/api/*", func(c *fiber.Ctx) error {
		// return serialized JSON.
		if c.Params("*") == "fruits" {

			response := FruitInterface{
				Page:   1,
				Fruits: []string{"apple", "peach", "pear", "watermelon"},
			}

			return c.JSON(response)

			// navigate to => http://localhost:3000/api/fruits
		}

		return c.SendString("API path: " + c.Params("*") + " -> do lookups with these values")
		// navigate to => http://localhost:3000/api/user/iggy
	})

	// Photon webhook testing.
	app.Post("/create-room", createRoom)
	app.Post("/close-room", closeRoom)

	// 404 handler.
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404) // => 404 "Not Found"
	})

	// Get port from env vars.
	var port = os.Getenv("PORT")

	// Use a default port if none was set in env.
	if port == "" {
		port = "3000"
	}

	// Start server on http://${heroku-url}:${port}
	app.Listen(":" + port)
}

// RoomResponse is what's send back to Photon.
type RoomResponse struct {
	State      string
	ResultCode int
}

func createRoom(c *fiber.Ctx) error {

	fmt.Println("Photon: create room: " + c.Request().String())

	var response = RoomResponse{
		"",
		0,
	}

	return c.JSON(response)
}

func closeRoom(c *fiber.Ctx) error {

	fmt.Println("Photon: close room:" + c.Request().String())

	var response = RoomResponse{
		"",
		0,
	}

	return c.JSON(response)
}
