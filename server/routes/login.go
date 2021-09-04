package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jakecoffman/crud"
	"net/http"
)

var oneYear = 60 * 60 * 24 * 256

type loginResponse struct {
	New bool
}

func login(c *gin.Context) {
	_, err := c.Cookie("player")
	if err == http.ErrNoCookie {
		playerId := uuid.New().String()
		// TODO are these values good?
		c.SetCookie("player", playerId, oneYear, "/", "", true, false)
		c.JSON(200, loginResponse{New: true})
		return
	}

	c.JSON(200, loginResponse{New: false})
}

var Login = crud.Spec{
	Method:      "GET",
	Path:        "/login",
	Handler:     login,
	Description: "",
	Tags:        []string{"Login"},
	Summary:     "Log in to the server",
	Validate:    crud.Validate{},
}
