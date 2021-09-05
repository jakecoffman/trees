package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jakecoffman/crud"
	"github.com/jakecoffman/trees/server/arcade"
)

var Admin = crud.Spec{
	Method:      "GET",
	Path:        "/count",
	Handler:     count,
	Description: "",
	Tags:        []string{"Login"},
	Summary:     "Checks ",
	Validate:    crud.Validate{},
}

type usersResponse struct {
	PlayerCount, GameCount int
}

func count(c *gin.Context) {
	playerCount, gameCount := arcade.Building.Counts()
	c.JSON(200, usersResponse{
		PlayerCount: playerCount,
		GameCount:   gameCount,
	})
}
