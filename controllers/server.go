package server

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"iggyzuk.com/go-server/models"
)

var roomCodes = make(map[string]*models.RoomCode) // A map of all allocated codes.
var freeCodes = make([]uint16, 9999)              // A slice of all the free codes.

// Run starts the photon server.
func Run() {
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

	var roomCode, err = getNextCode() // Generate a new code.

	// Error handling for when all codes are used up.
	if err != nil {
		return c.SendString(err.Error())
	}

	fmt.Println("Generated a new code: " + roomCode.Code)

	go timeoutRoom(roomCode)

	return c.SendString(roomCode.Code)
}

// createRoom is a handler for a Photon webhook that is called when a new room is created.
func createRoom(c *fiber.Ctx) error {

	fmt.Println("New room created.")

	// New room struct.
	room := new(models.CreateRoomRequest)

	// Parse body into struct.
	if err := c.BodyParser(room); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	if roomCode, ok := roomCodes[room.GameID]; ok {
		roomCode.Created = true // Confirm that the room was created by Photon.
		fmt.Println("Confirming that room: " + roomCode.Code + " has been created by Photon.")
		return c.JSON(models.RoomResponse{"", 0}) // Success.
	}

	// The code must have been removed – due to timeout.
	fmt.Println("No code for user: " + room.UserID)
	return c.Status(400).SendString("No code for user: " + room.UserID)
}

// closeRoom is a handler for a Photon webhook that is called when a room is closed
func closeRoom(c *fiber.Ctx) error {

	fmt.Println("Room closed.")

	// New room struct
	room := new(models.CloseRoomRequest)

	// Parse body into struct
	if err := c.BodyParser(room); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	returnCode(room.GameID)
	delete(roomCodes, room.GameID)

	fmt.Println("Room successfully removed: " + room.GameID)

	return c.JSON(models.RoomResponse{"", 0}) // Success.
}

// getNextCode generates the next unique room code.
func getNextCode() (*models.RoomCode, error) {

	if len(freeCodes) == 0 {
		return nil, errors.New("ran out of codes")
	}

	// Get a random free code.
	var freeCodeIndex = rand.Intn(len(freeCodes))
	var randomFreeCode = freeCodes[freeCodeIndex]

	// Remove (reslice) the free code from the list.
	freeCodes[freeCodeIndex] = freeCodes[len(freeCodes)-1]
	freeCodes = freeCodes[:len(freeCodes)-1]

	// Construct new room code object.
	var roomCode = new(models.RoomCode)
	roomCode.Code = fmt.Sprintf("%04d", int(randomFreeCode)) // Cast code from uint16 to string.
	roomCode.Created = false

	// Add it to the map.
	roomCodes[roomCode.Code] = roomCode

	// Return it.
	return roomCode, nil
}

// returnCode gives back the code to the be used again.
func returnCode(code string) {
	var codeInt, _ = strconv.Atoi(code)
	fmt.Println("Code returned: " + strconv.Itoa(codeInt))
	freeCodes = append(freeCodes, uint16(codeInt))
}

// timeoutRoom will remove a room if it was not created by Photon after 10 seconds.
func timeoutRoom(roomCode *models.RoomCode) {
	fmt.Println("Set timeout for room: " + roomCode.Code)
	time.Sleep(10 * time.Second)
	// If after a delay a room wasn't created by Photon we'll remove it.
	if !roomCode.Created {
		fmt.Println("Room: " + roomCode.Code + " timedout")
		returnCode(roomCode.Code)
		delete(roomCodes, roomCode.Code)
	} else {
		fmt.Println("Room: " + roomCode.Code + " was created, no need to timeout")
	}
}
