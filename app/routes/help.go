package routes

import (
	"main/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetHelp(c echo.Context) error {
	data := &models.APIHelp{
		Version:       "v2.0.0",
		WhatIsThis:    "API for this server, don't ask me how to use it... :(",
		ServerAddress: "play.tensyoserver.net",
		HomePage:      "https://www.tensyoserver.net",
		Wiki:          "https://wiki.tensyoserver.net/wiki/",
	}

	return c.JSON(http.StatusOK, data)
}
