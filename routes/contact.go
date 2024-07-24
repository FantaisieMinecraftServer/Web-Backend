package routes

import (
	"bytes"
	"encoding/json"
	"main/models"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/switchupcb/dasgo/dasgo"
)

func Create_contact(c echo.Context) error {
	name := c.FormValue("name")
	email := c.FormValue("email")
	objective := c.FormValue("objective")
	content := c.FormValue("content")

	webhook_data, _ := json.Marshal(
		dasgo.ExecuteWebhook{
			Username: dasgo.Pointer("お問い合わせ"),
			Embeds: []*dasgo.Embed{
				{
					Title: dasgo.Pointer("お問い合わせ - サイト"),
					Color: dasgo.Pointer(0x2b2d31),
					Fields: []*dasgo.EmbedField{
						{
							Name:   "名前",
							Value:  name,
							Inline: dasgo.Pointer(true),
						},
						{
							Name:   "メール",
							Value:  email,
							Inline: dasgo.Pointer(true),
						},
						{
							Name:   "メール",
							Value:  objective,
							Inline: dasgo.Pointer(true),
						},
						{
							Name:   "内容",
							Value:  content,
							Inline: dasgo.Pointer(false),
						},
					},
					Timestamp: dasgo.Pointer(dasgo.Timestamp(time.Now().Format(time.RFC3339))),
				},
			},
		},
	)

	res, err := http.Post(
		os.Getenv("WEBHOOK_URL"),
		"application/json",
		bytes.NewBuffer(webhook_data),
	)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Error{Reason: "unknown error"})
	}

	defer res.Body.Close()

	return c.Redirect(http.StatusFound, "https://www.tensyoserver.net/thanks")
}
