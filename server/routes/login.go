package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jakecoffman/crud"
	"github.com/jakecoffman/trees/server/arcade"
	"net/http"
	"net/url"
)

var Login = crud.Spec{
	Method:      "GET",
	Path:        "/login",
	Handler:     login,
	Description: "If you're already logged in this just returns your user information",
	Tags:        []string{"Login"},
	Summary:     "Returns login info, or logs you in anonymously",
	Validate:    crud.Validate{},
}

var oneYear = 60 * 60 * 24 * 256

type loginResponse struct {
	New  bool
	Code string
}

func login(c *gin.Context) {
	playerId, err := c.Cookie("player")
	if err == http.ErrNoCookie {
		playerId := uuid.New().String()
		// TODO are these values good?
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "player",
			Value:    url.QueryEscape(playerId),
			MaxAge:   oneYear,
			Path:     "/",
			Domain:   "trees.jakecoffman.com",
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
			HttpOnly: false,
		})
		c.JSON(200, loginResponse{New: true})
		return
	}

	response := loginResponse{New: false}
	player := arcade.Building.GetPlayer(playerId)
	if player != nil && player.Room != nil {
		response.Code = player.Room.Code
	}
	c.JSON(200, response)
}
