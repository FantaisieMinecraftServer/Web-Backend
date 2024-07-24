package models

import "github.com/switchupcb/dasgo/dasgo"

type Success struct {
	Status string `json:"status"`
}

type Webhook struct {
	Content   string        `json:"content"`
	Username  string        `json:"username"`
	AvatarUrl string        `json:"avatar_url"`
	Embeds    []dasgo.Embed `json:"embeds"`
}
