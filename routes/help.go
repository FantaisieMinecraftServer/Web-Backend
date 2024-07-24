package routes

import (
	"main/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Get_help(c echo.Context) error {
	return c.JSON(http.StatusOK, &models.APIHelp{
		Version:    "v2.1.0",
		WhatIsThis: "API for this server, don't ask me how to use it... :(",
		Author:     "https://github.com/FantaisieMinecraftServer",
		HomePage:   "https://www.tensyoserver.net",
	})
}
