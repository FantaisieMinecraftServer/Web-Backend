package models

type Success struct {
	Status string `json:"status"`
}

type Field struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type Embed struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Timestamp   string  `json:"timestamp"`
	Fields      []Field `json:"fields"`
}

type Webhook struct {
	Content   string  `json:"content"`
	Username  string  `json:"username"`
	AvatarUrl string  `json:"avatar_url"`
	Embeds    []Embed `json:"embeds"`
}
