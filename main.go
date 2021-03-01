package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Photon webhook testing.
	app.Post("/create-room", createRoom)
	app.Post("/close-room", closeRoom)
	app.Get("/room-code/*", roomCode)

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

var roomCodes = make(map[string]int)
var roomCount int

type createRoomRequest struct {
	ActorNr       int    `json:"ActorNr"`
	AppVersion    string `json:"AppVersion"`
	AppID         string `json:"AppId"`
	CreateOptions struct {
		MaxPlayers       int         `json:"MaxPlayers"`
		IsVisible        bool        `json:"IsVisible"`
		LobbyID          interface{} `json:"LobbyId"`
		LobbyType        int         `json:"LobbyType"`
		CustomProperties struct {
		} `json:"CustomProperties"`
		EmptyRoomTTL       int         `json:"EmptyRoomTTL"`
		PlayerTTL          int         `json:"PlayerTTL"`
		CheckUserOnJoin    bool        `json:"CheckUserOnJoin"`
		DeleteCacheOnLeave bool        `json:"DeleteCacheOnLeave"`
		SuppressRoomEvents bool        `json:"SuppressRoomEvents"`
		PublishUserID      bool        `json:"PublishUserId"`
		ExpectedUsers      interface{} `json:"ExpectedUsers"`
	} `json:"CreateOptions"`
	GameID   string `json:"GameId"`
	Region   string `json:"Region"`
	Type     string `json:"Type"`
	UserID   string `json:"UserId"`
	Nickname string `json:"Nickname"`
}

type closeRoomRequest struct {
	ActorCount int    `json:"ActorCount"`
	AppVersion string `json:"AppVersion"`
	AppID      string `json:"AppId"`
	GameID     string `json:"GameId"`
	Region     string `json:"Region"`
	Type       string `json:"Type"`
}

// RoomResponse is what's send back to Photon.
type RoomResponse struct {
	State      string
	ResultCode int
}

func createRoom(c *fiber.Ctx) error {
	// New room struct
	room := new(createRoomRequest)

	// Parse body into struct
	if err := c.BodyParser(room); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	// Generate code
	roomCount++
	roomCodes[room.GameID] = roomCount

	fmt.Println("Room Created:" + room.GameID + ", Code:" + strconv.Itoa(roomCodes[room.GameID]))
	fmt.Println("Details" + c.Request().String())

	var response = RoomResponse{
		"",
		0,
	}

	return c.JSON(response)
}

func closeRoom(c *fiber.Ctx) error {

	fmt.Println("Photon: close room:" + c.Request().String())

	// New room struct
	room := new(closeRoomRequest)

	// Parse body into struct
	if err := c.BodyParser(room); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	fmt.Println("Room Closed:" + room.GameID + ", Code:" + strconv.Itoa(roomCodes[room.GameID]))
	fmt.Println("Details" + c.Request().String())

	delete(roomCodes, room.GameID)

	var response = RoomResponse{
		"",
		0,
	}

	return c.JSON(response)
}

func roomCode(c *fiber.Ctx) error {

	var gameID = c.Params("*")

	if code, ok := roomCodes[gameID]; ok {
		var codeString = strconv.Itoa(code)

		fmt.Println("Get Code:" + gameID + ", Code:" + codeString)
		fmt.Println("Details" + c.Request().String())

		return c.SendString(codeString)
	}

	return c.SendStatus(404)
}
