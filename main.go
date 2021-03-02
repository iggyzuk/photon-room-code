package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// RoomCode is a server-side type that keeps track of all codes.
type RoomCode struct {
	code    string
	created bool
}

// CreateRoomRequest is a Photon type that is sent to us when a new room is created.
type CreateRoomRequest struct {
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

// CloseRoomRequest is a Photon type that is sent to us when a room is closed.
type CloseRoomRequest struct {
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

var roomCodes = make(map[string]*RoomCode) // A map of all allocated codes.
var freeCodes = make([]uint16, 5)          // A slide of all the free codes.

func main() {
	app := fiber.New()

	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator

	// Create free codes.
	for i := 0; i < len(freeCodes); i++ {
		freeCodes[i] = uint16(i)
	}

	// Photon webhook testing.
	app.Get("/room/gen_code", genCode)
	app.Post("/room/create", createRoom)
	app.Post("/room/close", closeRoom)

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

// genCode is a handler that is called by the user to get a new code.
func genCode(c *fiber.Ctx) error {

	var roomCode = getNextCode() // Generate a new code.

	fmt.Println("Generated a new code: " + roomCode.code)

	go timeoutRoom(roomCode)

	return c.SendString(roomCode.code)
}

// createRoom is a handler for a Photon webhook that is called when a new room is created.
func createRoom(c *fiber.Ctx) error {

	fmt.Println("New room created.")

	// New room struct.
	room := new(CreateRoomRequest)

	// Parse body into struct.
	if err := c.BodyParser(room); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	if roomCode, ok := roomCodes[room.GameID]; ok {
		roomCode.created = true // Confirm that the room was created by Photon.
		fmt.Println("Confirming that room: " + roomCode.code + " has been created by Photon.")
		return c.JSON(RoomResponse{"", 0}) // Success.
	}

	// The code must have been removed – due to timeout.
	fmt.Println("No code for user: " + room.UserID)
	return c.Status(400).SendString("No code for user: " + room.UserID)
}

// closeRoom is a handler for a Photon webhook that is called when a room is closed
func closeRoom(c *fiber.Ctx) error {

	fmt.Println("Room closed.")

	// New room struct
	room := new(CloseRoomRequest)

	// Parse body into struct
	if err := c.BodyParser(room); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	returnCode(room.GameID)
	delete(roomCodes, room.GameID)

	fmt.Println("Room successfully removed: " + room.GameID)

	return c.JSON(RoomResponse{"", 0}) // Success.
}

// getNextCode generates the next unique room code.
func getNextCode() *RoomCode {

	// Get a random free code.
	var freeCodeIndex = rand.Intn(len(freeCodes))
	var randomFreeCode = freeCodes[freeCodeIndex]

	// Remove (reslice) the free code from the list.
	freeCodes[randomFreeCode] = freeCodes[len(freeCodes)-1]
	freeCodes = freeCodes[:len(freeCodes)-1]

	// Construct new room code object.
	var roomCode = &RoomCode{
		strconv.Itoa(int(randomFreeCode)), // Cast code from uint16 to string.
		false,
	}

	// Add it to the map.
	roomCodes[roomCode.code] = roomCode

	// Return it.
	return roomCode
}

// returnCode gives back the code to the be used again.
func returnCode(code string) {
	var codeInt, _ = strconv.Atoi(code)
	freeCodes = append(freeCodes, uint16(codeInt))
}

// timeoutRoom will remove a room if it was not created by Photon after 10 seconds.
func timeoutRoom(roomCode *RoomCode) {
	fmt.Println("Set timeout for room: " + roomCode.code)
	time.Sleep(10 * time.Second)
	// If after a delay a room wasn't created by Photon we'll remove it.
	if !roomCode.created {
		fmt.Println("Room: " + roomCode.code + " timedout")
		returnCode(roomCode.code)
		delete(roomCodes, roomCode.code)
	} else {
		fmt.Println("Room: " + roomCode.code + " was created, no need to timeout")
	}
}
