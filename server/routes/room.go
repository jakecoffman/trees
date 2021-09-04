package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jakecoffman/crud"
	"github.com/jakecoffman/trees/server/arcade"
	"net/http"
)

var Room = []crud.Spec{
	{
		Method:      "POST",
		Path:        "/rooms",
		Handler:     createRoom,
		Description: "",
		Tags:        []string{"Rooms"},
		Summary:     "Create a game room",
		Validate:    crud.Validate{},
	},
	{
		Method:      "GET",
		Path:        "/rooms/:id",
		Handler:     getRoom,
		Description: "Mostly to see if the game is still running.",
		Tags:        []string{"Rooms"},
		Summary:     "Get a game room",
		Validate:    crud.Validate{},
	},
	{
		Method:      "DELETE",
		Path:        "/rooms/:id",
		Handler:     quitRoom,
		Description: "If you were playing in the game in the room, you will forfeit the game.",
		Tags:        []string{"Rooms"},
		Summary:     "Quit a game room.",
		Validate:    crud.Validate{},
	},
}

type roomResponse struct {
	Code string
}

func createRoom(c *gin.Context) {
	_, err := c.Cookie("player")
	if err != nil {
		c.AbortWithStatusJSON(401, "not logged in")
		return
	}
	room := arcade.NewRoom()
	c.JSON(http.StatusCreated, roomResponse{Code: room.Code})
}

func quitRoom(c *gin.Context) {
	playerId, err := c.Cookie("player")
	if err != nil {
		c.AbortWithStatusJSON(401, "not logged in")
		return
	}

	player := arcade.Building.GetPlayer(playerId)
	if player == nil || player.Room == nil {
		c.JSON(200, roomResponse{})
		return
	}
	// TODO player.Room shouldn't be a pointer
	room := player.Room
	room.Quit(player)
	arcade.Building.Shut(room)

	c.JSON(http.StatusOK, roomResponse{})
}

func getRoom(c *gin.Context) {
	_, err := c.Cookie("player")
	if err != nil {
		c.AbortWithStatusJSON(401, "not logged in")
		return
	}
	code := c.Param("id")
	room := arcade.Building.GetRoom(code)
	if room != nil {
		c.JSON(200, "OK")
	} else {
		c.JSON(404, "Not found")
	}
}
