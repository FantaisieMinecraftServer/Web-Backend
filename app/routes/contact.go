package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"main/models"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func PostContact(c echo.Context) error {
	name := c.FormValue("name")
	email := c.FormValue("email")
	objective := c.FormValue("objective")
	content := c.FormValue("content")

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	fields := []models.Field{{Name: "メールアドレス", Value: email, Inline: true}, {Name: "目的", Value: objective, Inline: true}}

	webhook_data, _ := json.Marshal(models.Webhook{Username: "サイト - お問い合わせ", Embeds: []models.Embed{{Title: "お問い合わせ - " + name, Description: content, Fields: fields, Timestamp: time.Now().UTC().Format("2006-01-02T15:04:05-0700")}}})

	res, _ := http.Post(
		os.Getenv("WEBHOOK_URL"),
		"application/json",
		bytes.NewBuffer(webhook_data),
	)

	body, _ := io.ReadAll(res.Body)

	fmt.Println(string(body))

	defer res.Body.Close()

	return c.Redirect(http.StatusFound, "https://www.tensyoserver.net/thanks")
}
